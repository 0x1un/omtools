package glimit

import "sync"

type pool struct {
	queue chan struct{}
	wg    *sync.WaitGroup
}

func New(size uint) *pool {
	if size == 0 {
		size = 1
	}
	return &pool{
		queue: make(chan struct{}, size),
		wg:    &sync.WaitGroup{},
	}
}

func (p *pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- struct{}{}
	}
	p.wg.Add(delta)
}

func (p *pool) Done() {
	<-p.queue
	p.wg.Done()
}

func(p *pool) Wait() {
	p.wg.Wait()
}