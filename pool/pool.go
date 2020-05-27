package pool

import (
	"log"
	"sync"
)

// Pool is a worker group that runs a number of jobs
type Pool struct {
	size   int
	jobs   chan func()
	wg     *sync.WaitGroup
	closed bool
}

// New pool
func New(size int) *Pool {
	return &Pool{
		size: size,
		jobs: make(chan func(), size*2),
		wg:   new(sync.WaitGroup),
	}
}

// Run pool workers
func (p *Pool) Run() {
	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go worker(i, p)
	}
}

// Wait until pool workers in processing, is a blocking operation
func (p *Pool) Wait() { p.wg.Wait() }

// Submit new job to the pool
func (p *Pool) Submit(job func()) {
	if p.closed {
		return
	}
	p.jobs <- job
}

// Close pool
func (p *Pool) Close() {
	log.Println("Shutdown worker pool")
	p.closed = true
	close(p.jobs)
}

func worker(id int, p *Pool) {
	defer p.wg.Done()

	log.Printf("Spawn worker [%d]", id)
	for job := range p.jobs {
		log.Printf("Worker [%d] received a job", id)
		job()
	}
	log.Printf("Shutdown worker [%d]", id)
}
