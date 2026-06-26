package metrics

import (
	"encoding/json"
	"net/http"
)

type QueueSizeReader interface {
	Size() int
}

type MetricsHandler struct {
	queue QueueSizeReader
}

func NewMetricsHandler(q QueueSizeReader) *MetricsHandler {
	return &MetricsHandler{queue: q}
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed: metrics requires HTTP GET"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GlobalMetrics.Snapshot(h.queue.Size()))
}
