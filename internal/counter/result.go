package counter

import "io"

type AsyncResult struct {
	reader io.ByteReader
	err    error
}

func (r *AsyncResult) Reader() io.ByteReader {
	return r.reader
}

func (r *AsyncResult) Error() error {
	return r.err
}

func NewAsyncResult(reader io.ByteReader, err error) AsyncResult {
	return AsyncResult{reader: reader, err: err}
}
