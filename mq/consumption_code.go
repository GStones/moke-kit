package mq

type ConsumptionCode int32

const (
	// "Success, let's keep going!"
	ConsumeAck ConsumptionCode = iota

	// "Success, but I'm done consuming!"
	ConsumeAckFinal

	// "Failure, but I'd like to try again! Send me that message again."
	ConsumeNackTransientFailure

	// "Failure, I give up on this message! Let's move on."
	ConsumeNackPersistentFailure
)
