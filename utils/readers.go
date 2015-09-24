package utils

import (
	"io"
	"io/ioutil"
	"log"
)

type multiReadCloser struct {
	readers []io.ReadCloser
	index   int
}

// MultiReadCloser returns a io.ReadCloser that's the logical concatenation of
// the provided input readers.
func MultiReadCloser(readerClosers ...io.ReadCloser) io.ReadCloser {

	return &multiReadCloser{
		readers: readerClosers,
	}
}

func (m *multiReadCloser) Read(p []byte) (int, error) {
	if m.index == len(m.readers) {
		return 0, io.EOF
	}

	size, err := m.readers[m.index].Read(p)
	if err != nil {
		if err != io.EOF {
			return size, err
		}
		if closeErr := m.readers[m.index].Close(); closeErr != nil {
			log.Printf("Got error while closing no longer needed readers inside multiReadCloser: %s\n", closeErr)
		}
		m.index++
		if m.index != len(m.readers) {
			err = nil
		}
	}

	return size, err
}

func (m *multiReadCloser) Close() error {
	c := new(CompositeError)
	for ; m.index < len(m.readers); m.index++ {
		err := m.readers[m.index].Close()
		if err != nil {
			c.AppendError(err)
		}
	}

	if c.Empty() {
		return nil
	}
	return c

}

type limitedReadCloser struct {
	io.ReadCloser
	maxLeft int
}

// LimitReadCloser wraps a io.ReadCloser but stops with EOF after `max` bytes.
func LimitReadCloser(readCloser io.ReadCloser, max int) io.ReadCloser {
	return &limitedReadCloser{
		ReadCloser: readCloser,
		maxLeft:    max,
	}
}

func (r *limitedReadCloser) Read(p []byte) (int, error) {
	readSize := min(r.maxLeft, len(p))
	size, err := r.ReadCloser.Read(p[:readSize])
	r.maxLeft -= size
	if r.maxLeft == 0 && err == nil {
		err = io.EOF
	}
	return size, err
}

func min(l, r int) int {
	if l > r {
		return r
	}
	return l
}

type skippingReadCloser struct {
	io.ReadCloser
	skip int64
}

// SkipReadCloser wraps a io.ReadCloser and ignores the first `skip` bytes.
func SkipReadCloser(readCloser io.ReadCloser, skip int64) io.ReadCloser {
	return &skippingReadCloser{
		ReadCloser: readCloser,
		skip:       skip,
	}
}

func (r *skippingReadCloser) Read(p []byte) (int, error) {
	if r.skip > 0 {
		if n, err := io.CopyN(ioutil.Discard, r.ReadCloser, r.skip); err != nil {
			return int(n), err
		}
		r.skip = 0
	}

	return r.ReadCloser.Read(p)
}
