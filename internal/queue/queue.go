package queue

import (
	"errors"

	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/ingestion"
)

var (
	ErrQueueFull = errors.New("pipeline capacity saturated: triggering active backpressure")
)

type MemoryQueue struct {
	channel chan ingestion.TelemetryEvent
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		channel: make(chan ingestion.TelemetryEvent, bufferSize),
	}
}

func (q *MemoryQueue) Push(event ingestion.TelemetryEvent) error {
	select {
	case q.channel <- event:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *MemoryQueue) Channel() <-chan ingestion.TelemetryEvent {
	return q.channel
}

func (q *MemoryQueue) Size() int {
	return len(q.channel)
}

func (q *MemoryQueue) Close() {
	close(q.channel)
}