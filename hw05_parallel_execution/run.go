package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type workerPool struct {
	wg       *sync.WaitGroup
	mx       *sync.Mutex
	stopped  bool
	errCount int
}

func newWorkerPool() *workerPool {
	return &workerPool{
		wg: new(sync.WaitGroup),
		mx: new(sync.Mutex),
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m > 0 {
		return runWithErrorCounting(tasks, n, m)
	}

	runIgnoringErrors(tasks, n)
	return nil
}

func runIgnoringErrors(tasks []Task, n int) {
	buf := make(chan struct{}, n)
	wg := &sync.WaitGroup{}

	for i := range tasks {
		wg.Add(1)
		buf <- struct{}{}
		go runTaskWithoutErrCh(tasks[i], wg, buf)
	}
	wg.Wait()
	close(buf)
}

func runWithErrorCounting(tasks []Task, n, m int) error {
	buf := make(chan struct{}, n)
	stopCh := make(chan struct{})
	wp := newWorkerPool()
	var err error
	for i := range tasks {
		if wp.stopped {
			err = ErrErrorsLimitExceeded
			break
		}
		wp.wg.Add(1)
		buf <- struct{}{}
		go wp.runTaskWithErrorCounting(tasks[i], buf, m)
	}

	wp.wg.Wait()
	close(stopCh)
	return err
}

func (wp *workerPool) runTaskWithErrorCounting(task Task, buf chan struct{}, maxCount int) {
	defer wp.wg.Done()
	defer wp.mx.Unlock()
	err := task()
	wp.mx.Lock()
	if err != nil {
		wp.errCount++
	}
	if wp.errCount >= maxCount {
		wp.stopped = true
		<-buf
		return
	}
	<-buf
}

func runTaskWithoutErrCh(task Task, wg *sync.WaitGroup, buf chan struct{}) {
	defer wg.Done()
	_ = task()
	<-buf
}
