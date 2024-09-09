package hw06pipelineexecution

import "sync/atomic"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	res := make(Bi)
	go func() {
		out := CreatePipeline(in, stages...)
		defer close(res)
		var stopped atomic.Bool

		for {
			select {
			case <-done:
				stopped.Store(true)
			case v, ok := <-out:
				if !ok {
					return
				}

				if !stopped.Load() {
					res <- v
				}
			}
		}
	}()

	return res
}

func CreatePipeline(in In, stages ...Stage) Out {
	if len(stages) == 0 {
		return nil
	}

	if len(stages) == 1 {
		return stages[0](in)
	}

	return CreatePipeline(stages[0](in), stages[1:]...)
}
