package reader

import (
	"bytes"
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
	reader := &bytes.Buffer{}

	for i := 0; i < bufferSize; i++ {
		reader.WriteByte(byte(i % 255))
	}

	assertions.NoError(buffer.write(reader))
	assertions.Equal(bufferSize, buffer.length)

	actual := make([]byte, 0, bufferSize)
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
