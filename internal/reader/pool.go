package reader

import (
	"runtime"
)

type pool struct {
	free chan *buffer
}

func newPool() *pool {
	maxProc := runtime.GOMAXPROCS(-1)
	p := &pool{
		free: make(chan *buffer, maxProc),
	}

	for i := 0; i < maxProc; i++ {
		p.free <- newFileBuf(p)
	}

	return p
}

func (p *pool) returnBuffer(b *buffer) {
	p.free <- b
}

func (p *pool) getBuffer() *buffer {
	return <-p.free
}
