package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

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
	errCh := make(chan error)
	stopCh := make(chan struct{})

	go func() {
		errCount := 0
		for range errCh {
			errCount++
			if errCount == m {
				stopCh <- struct{}{}
			}
		}
	}()

	wg := &sync.WaitGroup{}
	var err error
LOOP:
	for i := range tasks {
		select {
		case <-stopCh:
			err = ErrErrorsLimitExceeded
			break LOOP
		default:
			wg.Add(1)
			buf <- struct{}{}
			go runTaskWithErrCh(tasks[i], wg, buf, errCh)
		}
	}

	wg.Wait()
	close(buf)
	close(errCh)

	return err
}

func runTaskWithErrCh(task Task, wg *sync.WaitGroup, buf chan struct{}, errCh chan error) {
	defer wg.Done()
	err := task()
	<-buf
	if err != nil {
		errCh <- err
	}
}

func runTaskWithoutErrCh(task Task, wg *sync.WaitGroup, buf chan struct{}) {
	defer wg.Done()
	_ = task()
	<-buf
}
