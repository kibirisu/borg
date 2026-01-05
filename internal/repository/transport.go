package repository

type Queue chan Job

type Job struct{}

type TransportRepository interface{}

type transportRepository struct {
	queue Queue
}

var _ TransportRepository = (*transportRepository)(nil)

func NewTransport() TransportRepository {
	ch := make(Queue, 5)
	repo := transportRepository{ch}
	return &repo
}

func foo() {
	// ch := make(Queue, 0)
}
