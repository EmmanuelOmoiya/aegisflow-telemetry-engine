package worker

import (
	"log"
	"sync"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/ingestion"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/metrics"
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

func (w *Worker) Start() {
	w.waitGroup.Add(1)

	go func() {
		defer w.waitGroup.Done()
		
		log.Printf("[Worker #%d] Execution pipeline spawned successfully. Awaiting telemetry logs...\n", w.id)
		for event := range w.queue.Channel() {
			// Phase 1 Placeholder: Log the incoming target package trace for local auditing
			log.Printf("[Worker #%d] Processing trace - Device: %s, Source: %s, Event: %s\n", 
				w.id, event.DeviceID, event.Source, event.EventType,
			)
			metrics.GlobalMetrics.IncrementProcessed()
		}

		log.Printf("[Worker #%d] Ingestion channel closed. Cleaning up routing allocations and exiting safely.\n", w.id)
	}()
}