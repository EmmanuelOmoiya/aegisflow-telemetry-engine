package worker

import (
	"log"
	"sync"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/evaluator"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/queue"
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
			log.Printf("[Worker #%d] Processing trace - Device: %s, Source: %s, Event: %s\n", 
				w.id, event.DeviceID, event.Source, event.EventType,
			)

			metrics.GlobalMetrics.IncrementProcessed()

			evalContext := map[string]interface{}{
				"device_id":  event.DeviceID,
				"source":     event.Source,
				"event_type": event.EventType,
				"payload":    event.Payload, // Supports nested structural field resolution (e.g. payload.temperature)
			}

			if rule, exists := evaluator.GlobalRegistry.Lookup("critical_anomaly_rule"); exists {
				result, err := rule.Evaluate(evalContext)
				if err != nil {
					log.Printf("[Worker #%d Error] Rule validation evaluation anomaly: %v\n", w.id, err)
					continue
				}

				if isViolated, ok := result.(bool); ok && isViolated {
					log.Printf("[ALERT] Policy violation detected on device %s via rule engine constraint analysis!\n", event.DeviceID)
					metrics.GlobalMetrics.IncrementViolation()
				}
			}
		}

		log.Printf("[Worker #%d] Ingestion channel closed. Cleaning up routing allocations and exiting safely.\n", w.id)
	}()
}