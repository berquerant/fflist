package iox

import (
	"errors"
	"io"
)

type ReaderAndCloser interface {
	Reader() io.Reader
	Close() error
}

func NewMultiReaderAndCloser(rs ...io.ReadCloser) *MultiReaderAndCloser {
	return &MultiReaderAndCloser{
		rs: rs,
	}
}

type MultiReaderAndCloser struct {
	rs []io.ReadCloser
}

func (r *MultiReaderAndCloser) Reader() io.Reader {
	rs := make([]io.Reader, len(r.rs))
	for i, r := range r.rs {
		rs[i] = r
	}
	return io.MultiReader(rs...)
}

func (r *MultiReaderAndCloser) Close() error {
	var errs []error
	for _, x := range r.rs {
		if err := x.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type readerAndCloser struct {
	r io.Reader
}

func (r *readerAndCloser) Reader() io.Reader { return r.r }
func (readerAndCloser) Close() error         { return nil }

func AsReaderAndCloser(r io.Reader) ReaderAndCloser {
	return &readerAndCloser{r}
}
