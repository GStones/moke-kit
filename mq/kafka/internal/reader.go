package internal

import (
	"context"
	"sync"

	k "github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/mq"
	"github.com/gstones/platform/services/common/utils"
)

type reader struct {
	logger           *StatsLogger
	mutex            sync.Mutex
	kReader          *k.Reader
	returnPool       chan *reader
	deadLetterWriter *deadLetterWriter

	decoder     mq.Decoder
	vPtrFactory mq.ValuePtrFactory
	isReturning bool
	cancel      context.CancelFunc
}

const (
	maxCommitAttempts = 3
)

func NewReader(
	logger *zap.Logger,
	config k.ReaderConfig,
	pool chan *reader,
	statsPeriod StatsPeriod,
	deadLetterWriter *deadLetterWriter,
) *reader {
	kReader := k.NewReader(config)

	statsLogger := NewStatsLogger(
		logger,
		statsPeriod,
		&ReaderStatsProvider{reader: kReader},
	)

	statsLogger.Start()

	return &reader{
		kReader:          kReader,
		returnPool:       pool,
		logger:           statsLogger,
		deadLetterWriter: deadLetterWriter,
	}
}

func (r *reader) ReceiveLoop(
	handler mq.SubResponseHandler,
	decoder mq.Decoder,
	vPtrFactory mq.ValuePtrFactory,
) {
	go r.receiveLoop(handler, decoder, vPtrFactory)
}

// receiveLoop() does the following:
// (1) Fetches messages from the kafka instance
// (2) Delivers messages to the provided handler
// (3) Takes actions based the handler's response to include:
//
//	Move on to next msg; Re-attempt msg delivery; Send msg to deadLetter queue; Quit out of msg fetch/delivery loop
func (r *reader) receiveLoop(
	handler mq.SubResponseHandler,
	decoder mq.Decoder,
	vPtrFactory mq.ValuePtrFactory,
) {
	ctx := r.initialization(handler, decoder, vPtrFactory)

FetchLoop:
	for {
		if r.IsReturning() {
			break FetchLoop
		}

		if message, kMessage, err := r.fetchNextMessage(ctx); err != nil {
			handler(nil, err)
			break FetchLoop
		} else {

		DeliveryRetryLoop:
			for {
				if r.IsReturning() {
					break FetchLoop
				}

				// Deliver next kafka message to the provided subscription handler.
				clientRespCode := handler(message, nil)

				// Handle the handler's response to delivery attempt.
				// Depending on the returned handler response code, we may take various actions.
				switch clientRespCode {
				// "Success, let's keep going!"
				case mq.ConsumeAck:
					if err := r.commitMessage(kMessage); err != nil {
						handler(nil, err)
						break FetchLoop
					} else {
						continue FetchLoop
					}

				// "Success, but I'm done consuming!"
				case mq.ConsumeAckFinal:
					handler(nil, r.commitMessage(kMessage))

				// "Failure, but I'd like to try again! Send me that message again."
				case mq.ConsumeNackTransientFailure:
					continue DeliveryRetryLoop

				// "Failure, I give up on this message! Let's move on."
				case mq.ConsumeNackPersistentFailure:
					if err := r.commitMessage(kMessage); err != nil {
						handler(nil, err)
						break FetchLoop
					} else {
						r.deadLetterWriter.Write(kMessage.Value)
						continue FetchLoop
					}

				// Unsupported ConsumptionCode.
				default:
					handler(nil, mq.ErrUnsupportedConsumeCode)
					if err := r.commitMessage(kMessage); err != nil {
						break FetchLoop
					} else {
						r.deadLetterWriter.Write(kMessage.Value)
						break FetchLoop
					}
				}
			}
		}
	}

	r.teardown()
	return
}

func (r *reader) ReturnToPool() {
	r.mutex.Lock()

	r.isReturning = true

	// For the case where Unsubscribe() is called immediately upon Subscription creation
	if r.cancel != nil {
		r.cancel()
	}

	r.mutex.Unlock()
}

func (r *reader) fetchNextMessage(ctx context.Context) (mq.Message, k.Message, error) {
	if kMessage, err := r.kReader.FetchMessage(ctx); err != nil {
		return nil, k.Message{}, err
	} else if uuid, err := utils.NewUUIDString(); err != nil {
		return nil, k.Message{}, err
	} else {
		var message mq.Message

		if r.decoder != nil && r.vPtrFactory != nil {
			vPtrMessage := r.vPtrFactory.NewVPtr()

			if err := r.decoder.Decode(kMessage.Topic, kMessage.Value, vPtrMessage); err != nil {
				return nil, k.Message{}, err
			} else {
				message = NewMessage(uuid, kMessage.Topic, nil, vPtrMessage)
			}
		} else {
			message = NewMessage(uuid, kMessage.Topic, kMessage.Value, nil)
		}

		return message, kMessage, nil
	}
}

func (r *reader) commitMessage(msg k.Message) (err error) {
	for i := 0; i < maxCommitAttempts; i++ {
		if err = r.kReader.CommitMessages(context.Background(), msg); err != nil {
			continue
		} else {
			return nil
		}
	}

	return err
}

func (r *reader) initialization(
	handler mq.SubResponseHandler,
	decoder mq.Decoder,
	vPtrFactory mq.ValuePtrFactory,
) context.Context {
	r.mutex.Lock()

	r.decoder = decoder
	r.vPtrFactory = vPtrFactory
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	r.mutex.Unlock()

	return ctx
}

// Upon returning to pool, the reader sets appropriate decoder default values
func (r *reader) teardown() {
	r.mutex.Lock()

	r.decoder = nil
	r.vPtrFactory = nil
	r.isReturning = false
	r.returnPool <- r

	r.mutex.Unlock()
}

// Use IsReturning() from within a mutex lock and isReturning from outside a lock
func (r *reader) IsReturning() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.isReturning
}
