package interval

import "fmt"

type Tasker interface {
	Handle(task any)
}

const (
	// RR select que by using Round Robin
	RR = 0
	// SRC select que by using outer source
	SRC = 1
)

type WorkerPool struct {
	MaxPoolSize int
	MaxQueLen   int
	taskQueue   []chan any
	tasker      Tasker
	qid         int
	Mod         int
}

func NewWorkPool(maxPoolSize, maxQueLen int, tasker Tasker, mod int) *WorkerPool {
	return &WorkerPool{
		MaxPoolSize: maxPoolSize,
		MaxQueLen:   maxQueLen,
		taskQueue:   make([]chan any, maxPoolSize),
		tasker:      tasker,
		qid:         0,
		Mod:         mod,
	}
}

func (w *WorkerPool) work(qid int) {
	for task := range w.taskQueue[qid] {
		w.tasker.Handle(task)
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.MaxPoolSize; i++ {
		w.taskQueue[i] = make(chan any, w.MaxQueLen)
		go w.work(i)
	}
}

// AppendTask used Round Robin or Source
func (w *WorkerPool) AppendTask(task any, src int) {
	switch w.Mod {
	case RR:
		w.qid = w.qid % w.MaxPoolSize
		w.taskQueue[w.qid] <- task
		fmt.Printf("task que %d recv a task\n", w.qid)
		w.qid++
	case SRC:
		qid := src % w.MaxPoolSize
		w.taskQueue[qid] <- task
		fmt.Printf("task que %d recv a task\n", w.qid)
	}

}

func (w *WorkerPool) Shut() {
	for i := 0; i < w.MaxPoolSize; i++ {
		close(w.taskQueue[i])
	}
}
