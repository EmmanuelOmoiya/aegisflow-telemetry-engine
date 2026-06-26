package worker

import (
	"log"
	"sync"
)

type WorkerPool struct {
	workerCount int
	queue       IngestionQueueReader
	workers     []*Worker
	waitGroup   sync.WaitGroup
}

func NewWorkerPool(workerCount int, queue IngestionQueueReader) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		queue:       queue,
		workers:     make([]*Worker, 0, workerCount),
	}
}

func (p *WorkerPool) StartPool() {
	log.Printf("[WorkerPool] Igniting processing core cluster. Spawning %d background workers...\n", p.workerCount)

	for i := 1; i <= p.workerCount; i++ {
		wrk := NewWorker(i, p.queue, &p.waitGroup)
		
		p.workers = append(p.workers, wrk)
		
		wrk.Start()
	}

	log.Printf("[WorkerPool] All %d background workers successfully online and routing.\n", p.workerCount)
}