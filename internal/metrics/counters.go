package metrics

import (
	"sync/atomic"
)

type OperationalMetrics struct {
	IngestedCount  uint64
	ProcessedCount uint64
	ViolationCount uint64
}

var GlobalMetrics = &OperationalMetrics{}

func (m *OperationalMetrics) IncrementIngested() {
	atomic.AddUint64(&m.IngestedCount, 1)
}

func (m *OperationalMetrics) IncrementProcessed() {
	atomic.AddUint64(&m.ProcessedCount, 1)
}

func (m *OperationalMetrics) IncrementViolated() {
	atomic.AddUint64(&m.ViolationCount, 1)
}

func (m *OperationalMetrics) Snapshot(currentQueueSize int) map[string]uint64 {
	return map[string]uint64{
		"ingested_events":   atomic.LoadUint64(&m.IngestedCount),
		"processed_events":  atomic.LoadUint64(&m.ProcessedCount),
		"policy_violations": atomic.LoadUint64(&m.ViolationCount),
		"queue_depth_items": uint64(currentQueueSize),
	}
}
