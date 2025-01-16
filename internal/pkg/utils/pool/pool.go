package pool

import (
	"context"
)

type Task[I, O any] func(I) O

type Pool[I, O any] struct {
	ch   chan struct{}
	size int
}

func NewPool[I, O any](size int) *Pool[I, O] {
	return &Pool[I, O]{
		ch:   make(chan struct{}, size),
		size: size,
	}
}

func (p *Pool[I, O]) Submit(ctx context.Context, task Task[I, O], input I) <-chan O {
	resultChan := make(chan O, 1)

	go func() {
		select {
		case p.ch <- struct{}{}:
			defer func() {
				<-p.ch
				close(resultChan)
			}()

			result := task(input)
			select {
			case <-ctx.Done():
				return
			case resultChan <- result:
			}

		case <-ctx.Done():
			close(resultChan)
			return
		}
	}()

	return resultChan
}
