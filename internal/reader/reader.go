package reader

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/stainour/test12/internal/counter"
)

type fileReader struct {
	walker  fileWalker
	pool    *pool
	once    sync.Once
	results chan counter.AsyncResult
}

func NewFileReader(path string) (counter.AsyncByteReader, error) {
	walker, err := newFileWalker(path)
	if err != nil {
		return nil, fmt.Errorf("error creating file walker: %w", err)
	}

	return &fileReader{
		walker:  walker,
		pool:    newPool(),
		results: make(chan counter.AsyncResult, runtime.GOMAXPROCS(-1)),
		once:    sync.Once{},
	}, nil
}

func (r *fileReader) Readers(ctx context.Context) <-chan counter.AsyncResult {
	r.once.Do(func() {
		filePaths := r.walker.walk(ctx)
		go func() {
			r.readFiles(ctx, filePaths)
			close(r.results)
		}()
	})

	return r.results
}

func (r *fileReader) readFiles(ctx context.Context, filePaths chan walkResult) {
	for p := range filePaths {
		if p.err != nil {
			r.results <- counter.NewAsyncResult(nil, fmt.Errorf("error reading path file: %w", p.err))
			return
		}

		f, err := os.Open(p.path)
		if err != nil {
			r.results <- counter.NewAsyncResult(nil, fmt.Errorf("error opening file %s: %w", p.path, err))
			return
		}

		for {
			select {
			case <-ctx.Done():
				return

			default:
			}

			b := r.pool.getBuffer()

			n, err := f.Read(b.buf)
			if err == io.EOF {
				r.pool.returnBuffer(b)
				break
			} else if err != nil {
				r.results <- counter.NewAsyncResult(nil, fmt.Errorf("error reading file %s: %w", p.path, err))
				return
			}

			b.updateLength(n)
			r.results <- counter.NewAsyncResult(b, nil)
		}

		err = f.Close()
		if err != nil {
			r.results <- counter.NewAsyncResult(nil, fmt.Errorf("error closing file %s: %w", p.path, err))
			return
		}
	}
}
