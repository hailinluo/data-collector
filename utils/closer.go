package utils

import (
	"github.com/hashicorp/go-multierror"
	"io"
)

type CloseFunc func() error

func (c CloseFunc) Close() error {
	return c()
}

type Closer struct {
	Closers []io.Closer
}

func (c *Closer) AppendCloser(sc io.Closer) {
	c.Closers = append(c.Closers, sc)
}

func (c *Closer) AppendCloseFunc(f func() error) {
	c.Closers = append(c.Closers, CloseFunc(f))
}

func (c *Closer) Close() error {
	var err error
	var mErr = new(multierror.Error)
	var closerLen = len(c.Closers)
	for i := closerLen - 1; i >= 0; i-- {
		err = c.Closers[i].Close()
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}
	return mErr.ErrorOrNil()
}
