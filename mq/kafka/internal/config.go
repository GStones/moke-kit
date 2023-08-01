package internal

import (
	"time"

	"github.com/gstones/platform/services/common/mq"
)

// TODO: Per Steve, configuration should be hardcoded (no env vars)
// TODO: until we have configured Kafka to create topics on demand.
// TODO: See GitLab issue #602

const (
	numPartitions               = 16
	replicationFactor           = 1
	deadLetterNumPartitions     = 16
	deadLetterReplicationFactor = 1
	readerStatsPeriod           = time.Duration(1 * time.Minute)
	writerStatsPeriod           = time.Duration(1 * time.Minute)
	balancer                    = BalancerCodeRoundRobin
	defaultDeliverySemantics    = mq.AtMostOnce
)
