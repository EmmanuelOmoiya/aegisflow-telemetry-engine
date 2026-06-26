package worker

import (
	"log"
	"sync"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/ingestion"
)

type IngestionQueueReader interface {
	Channel() <-chan ingestion.TelemetryEvent
}

type Worker struct {
	id         int
	queue      IngestionQueueReader
	quitChan   chan struct{}
	waitGroup  *sync.WaitGroup
}

func NewWorker(id int, queue IngestionQueueReader, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:         id,
		queue:      queue,
		quitChan:   make(chan struct{}),
		waitGroup:  wg,
	}
}