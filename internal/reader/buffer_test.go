package reader

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	assertions := assert.New(t)

	pool := newPool()
	poolSize := len(pool.free)

	buffer := pool.getBuffer()
	bufferSize := len(buffer.buf)
	actual := make([]byte, 0, bufferSize)

	for i := 0; i < bufferSize; i++ {
		buffer.buf[i] = byte(i % 255)
	}
	buffer.updateLength(bufferSize)
	assertions.Equal(bufferSize, buffer.length)

	for {
		b, err := buffer.ReadByte()
		if err == io.EOF {
			break
		} else {
			assertions.NoError(err)
		}

		actual = append(actual, b)
	}

	assertions.Equal(0, buffer.position)
	assertions.Equal(0, buffer.length)
	assertions.Equal(buffer.buf, actual)
	assertions.True(buffer.freed)
	assertions.Equal(poolSize, len(pool.free))
}
