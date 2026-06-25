package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/config"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/ingestion"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/metrics"
	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/queue"
)

func main() {
	log.Println("Initializing AegisFlow Telemetry Engine core setup...")

	cfg := config.LoadConfig()

	memQueue := queue.NewMemoryQueue(cfg.QueueBuffer)
	log.Printf("[Queue] Bounded pipeline allocated with max capacity buffer: %d\n", cfg.QueueBuffer)

	ingestHandler := ingestion.NewIngestionHandler(memQueue)
	metricsHandler := metrics.NewMetricsHandler()

	mux := http.NewServeMux()
	mux.Handle("/v1/telemetry", ingestHandler)
	mux.Handle("/metrics", metricsHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("[Server] Ingestion gateway online and listening on channel port %s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[Fatal] Server crashed unexpectedly: %v\n", err)
		}
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdownChan
	log.Printf("[Shutdown] Signal received: %v. Cleaning up ingestion channels cleanly...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[Shutdown] Force closing active instances: %v\n", err)
	}

	log.Println("AegisFlow server connection pool closed down with zero leak footprint. Exit success.")
}
