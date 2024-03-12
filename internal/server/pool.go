package server

import (
	"errors"
	"runtime"
	"sync"
)

var (
	errWorker = errors.New("all workers is busying")
	errCap    = errors.New("initCap must <= maxCap")
	errArgs   = errors.New("invalied args num, need most two args")
)

var defaultCap = uint(runtime.NumCPU() * 100)

type WorkerPool interface {
	Get() (Worker, error)
	Put(Worker) error
}

type defaultWorkerPool struct {
	initCap uint
	maxCap  uint
	mu      sync.Mutex
	workers []Worker
	used    uint
}

func (wp *defaultWorkerPool) Get() (Worker, error) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.used >= wp.maxCap {
		wp.used = wp.maxCap
		return nil, errWorker
	}

	if wp.workers[wp.used] == nil {
		defaultWorker := newWorker(wp.used)
		defaultWorker.run()
		wp.workers[wp.used] = defaultWorker
	}

	worker := wp.workers[wp.used]
	wp.used++
	return worker, nil
}

func (wp *defaultWorkerPool) Put(worker Worker) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	w := worker.(*defaultWorker)
	wp.used--
	wp.workers[w.pos] = wp.workers[wp.used]
	wp.workers[w.pos].(*defaultWorker).pos = w.pos
	w.pos = wp.used
	wp.workers[wp.used] = w
	return nil
}

func (wp *defaultWorkerPool) init() (WorkerPool, error) {
	wp.workers = make([]Worker, wp.maxCap)
	for i := uint(0); i < wp.initCap; i++ {
		wp.workers[i] = newWorker(i)
	}
	return wp, nil
}

func (wp *defaultWorkerPool) initWorkers() {
	for i := uint(0); i < wp.initCap; i++ {
		wp.workers[i].(*defaultWorker).run()
	}
}

func DefaultWorkerPool(caps ...uint) (WorkerPool, error) {
	wp := &defaultWorkerPool{}
	switch len(caps) {
	case 0:
		wp.initCap = defaultCap
		wp.maxCap = defaultCap
	case 1:
		wp.initCap = uint(caps[0])
		wp.maxCap = wp.initCap
	case 2:
		wp.initCap = uint(caps[0])
		wp.maxCap = uint(caps[1])
		if wp.maxCap < wp.initCap {
			return nil, errCap
		}
	default:
		return nil, errArgs
	}
	wp.init()
	wp.initWorkers()
	return wp, nil
}
