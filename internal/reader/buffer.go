package reader

import (
	"io"
)

const bufSize = 2 * 1024 * 1024

type buffer struct {
	buf      []byte
	position int
	length   int
	pool     *pool
	freed    bool
}

func newFileBuf(p *pool) *buffer {
	return &buffer{
		buf:  make([]byte, bufSize),
		pool: p,
	}
}

func (b *buffer) reset() {
	b.position = 0
	b.length = 0
}

func (b *buffer) ReadByte() (byte, error) {
	if b.position >= b.length {
		b.reset()

		if !b.freed {
			b.freed = true
			b.pool.returnBuffer(b)
		}

		return 0, io.EOF
	}

	v := b.buf[b.position]
	b.position++

	return v, nil
}

func (b *buffer) updateLength(length int) {
	b.length = length
	b.freed = false
}
