package ingestion

import (
	"encoding/json"
	"net/http"

	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/metrics"
)

type IngestionQueue interface {
	Push(event TelemetryEvent) error
}

type IngestionHandler struct {
	queue IngestionQueue
}

func NewIngestionHandler(q IngestionQueue) *IngestionHandler {
	return &IngestionHandler{queue: q}
}

func (h *IngestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed: ingestion requires HTTP POST"})
		return
	}

	var event TelemetryEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json payload structure"})
		return
	}

	if err := event.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	metrics.GlobalMetrics.IncrementIngested()

	if err := h.queue.Push(event); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "rejected", "message": "server telemetry ingress pipeline is completely saturated"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted", "message": "telemetry packet queued successfully for background worker evaluation"})
}
