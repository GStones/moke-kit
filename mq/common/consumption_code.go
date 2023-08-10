package common

type ConsumptionCode int32

const (
	// ConsumeAck "Success, let's keep going!"
	ConsumeAck ConsumptionCode = iota

	// ConsumeAckFinal "Success, but I'm done consuming!"
	ConsumeAckFinal

	// ConsumeNackTransientFailure "Failure, but I'd like to try again! Send me that message again."
	ConsumeNackTransientFailure

	// ConsumeNackPersistentFailure "Failure, I give up on this message! Let's move on."
	ConsumeNackPersistentFailure
)
