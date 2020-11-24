package counter

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"
)

type AsyncByteReader interface {
	// Should be thread safe
	Readers(ctx context.Context) <-chan AsyncResult
}

const maxASCIICode = 128

// AsciiCodeCount map ASCII code into symbol frequency count.
type ASCIICodeCount [maxASCIICode]int32

type ASCIICounter struct {
	reader AsyncByteReader
}

func NewASCIICounter(r AsyncByteReader) *ASCIICounter {
	return &ASCIICounter{
		reader: r,
	}
}

func (c *ASCIICounter) Count(ctx context.Context) (ASCIICodeCount, error) {
	maxProc := runtime.GOMAXPROCS(-1)
	errors := make(chan error, maxProc)
	sum := make(chan *ASCIICodeCount, maxProc)

	ctx, cancel := context.WithCancel(ctx)

	group := sync.WaitGroup{}
	group.Add(maxProc)

	for i := 0; i < maxProc; i++ {
		go func() {
			err := c.countReaders(ctx, cancel, sum)
			if err != nil {
				errors <- err
			}

			group.Done()
		}()
	}

	result := make(chan ASCIICodeCount)

	go func() {
		counts := ASCIICodeCount{}

		for c := range sum {
			for i := range c {
				counts[i] += c[i]
			}
		}
		result <- counts
	}()

	group.Wait()
	close(errors)
	close(sum)

	for e := range errors {
		return ASCIICodeCount{}, e
	}

	return <-result, nil
}

func (c *ASCIICounter) countReaders(ctx context.Context, cancel context.CancelFunc, sum chan<- *ASCIICodeCount) error {
	counts := &ASCIICodeCount{}
	readers := c.reader.Readers(ctx)

	for {
		select {
		case r, ok := <-readers:
			if !ok {
				sum <- counts
				return nil
			}

			if r.Error() != nil {
				return r.Error()
			}

			br := r.Reader()

			for {
				b, err := br.ReadByte()
				if err == io.EOF {
					break
				} else if err != nil {
					return fmt.Errorf("error reading byte: %w", err)
				}

				if b >= maxASCIICode {
					cancel()
					return fmt.Errorf("invalid ASCII symbol code %x", b)
				}
				counts[b]++
			}

			break

		case <-ctx.Done():
			return nil
		}
	}
}
