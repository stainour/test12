package counter_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stainour/test12/internal/counter"
)

type mockReader struct {
	r chan counter.AsyncResult
}

func newMockReader() *mockReader {
	return &mockReader{r: make(chan counter.AsyncResult, 1000)}
}

func (r *mockReader) close() {
	close(r.r)
}

func (r *mockReader) Readers(ctx context.Context) <-chan counter.AsyncResult {
	return r.r
}

func TestCountErrorReader(t *testing.T) {
	assertions, counter := arrange(t, counter.NewAsyncResult(nil, errors.New("error")))

	_, err := counter.Count(context.Background())
	assertions.Error(err)
}

func TestCountInvalidASCIICode(t *testing.T) {
	assertions, counter := arrange(t, counter.NewAsyncResult(bytes.NewBuffer([]byte{1, 5, 6, 8, 90, 150}), nil))

	_, err := counter.Count(context.Background())
	assertions.Error(err)
}

func TestCount(t *testing.T) {
	const bufSize = 128

	var results []counter.AsyncResult

	for i := 0; i < bufSize; i++ {
		buffer := &bytes.Buffer{}
		buffer.WriteByte(byte(i))
		results = append(results, counter.NewAsyncResult(buffer, nil))
	}

	expected := counter.ASCIICodeCount{}

	for i := 0; i < bufSize; i++ {
		expected[i] = 1
	}

	assertions, counter := arrange(t, results...)

	counts, err := counter.Count(context.Background())

	assertions.NoError(err)
	assertions.Equal(expected, counts)
}

func arrange(t *testing.T, r ...counter.AsyncResult) (*assert.Assertions, *counter.ASCIICounter) {
	assertions := assert.New(t)
	reader := newMockReader()
	counter := counter.NewASCIICounter(reader)

	for _, v := range r {
		reader.r <- v
	}

	reader.close()

	return assertions, counter
}
