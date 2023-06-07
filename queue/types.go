package queue

type Queue interface {
	Enqueue(val any) error
	Dequeue() (val any, err error)
}
