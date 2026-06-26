package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/EmmanuelOmoiya/aegisflow/internal/ingestion"
	"github.com/EmmanuelOmoiya/aegisflow/internal/metrics"
	"github.com/EmmanuelOmoiya/aegisflow/internal/queue"
	"github.com/EmmanuelOmoiya/aegisflow/internal/worker"
)

func TestHighVelocityConcurrentPipeline(t *testing.T) {
	metrics.GlobalMetrics.IngestedCount = 0
	metrics.GlobalMetrics.ProcessedCount = 0
	metrics.GlobalMetrics.ViolationCount = 0

	concurrentRequests := 1000
	memQueue := queue.NewMemoryQueue(concurrentRequests)

	workerCount := 10
	workerPool := worker.NewWorkerPool(workerCount, memQueue)
	workerPool.StartPool()

	handler := ingestion.NewIngestionHandler(memQueue)
	server := httptest.NewServer(handler)
	defer server.Close()

	var wg sync.WaitGroup
	payload := []byte(`{
		"device_id": "device-test-xyz",
		"source": "sensor_array_alpha",
		"event_type": "temperature_drift",
		"metadata": {"reading": 42.7}
	}`)

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				t.Errorf("HTTP request failed completely: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusAccepted {
				t.Errorf("Expected status code 202 Accepted, got: %d", resp.StatusCode)
			}
		}()
	}

	wg.Wait()

	memQueue.Close()
	workerPool.Stop()

	snapshot := metrics.GlobalMetrics.Snapshot(memQueue.Size())

	if snapshot["ingested_events"] != uint64(concurrentRequests) {
		t.Errorf("Data leakage detected! Expected %d ingested events, but metrics recorded: %d", concurrentRequests, snapshot["ingested_events"])
	}

	if snapshot["processed_events"] != uint64(concurrentRequests) {
		t.Errorf("Processing dropped! Expected %d processed events, but background workers only finished: %d", concurrentRequests, snapshot["processed_events"])
	}

	if snapshot["queue_depth_items"] != 0 {
		t.Errorf("Memory drain failure! Expected remaining queue depth to be 0, but got: %d", snapshot["queue_depth_items"])
	}
}