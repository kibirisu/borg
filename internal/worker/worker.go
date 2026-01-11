package worker

import (
	"context"
	"log"
)

type Job func(context.Context) error

type Worker interface {
	Enqueue(Job)
	Cancel()
}

type worker struct {
	queue chan Job
}

var _ Worker = (*worker)(nil)

func New(ctx context.Context) Worker {
	w := &worker{make(chan Job)}
	for range 5 {
		go w.spawn(ctx)
	}
	return w
}

func (w *worker) spawn(ctx context.Context) {
	for job := range w.queue {
		if err := job(ctx); err != nil {
			log.Println(err)
		}
	}
}

// Enqueue implements Worker.
func (w *worker) Enqueue(job Job) {
	w.queue <- job
}

// Cancel implements Worker.
func (w *worker) Cancel() {
	close(w.queue)
}
