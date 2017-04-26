package errchan

import (
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

func formatFile(f string) string {
	fileName := filepath.Base(f)
	dir := filepath.Dir(f)
	return filepath.Join(filepath.Base(dir), fileName)
}

type ErrorProducer interface {
	Errors() <-chan error
}

type ErrChan struct {
	errors chan error
}

func New(b int) *ErrChan {
	return &ErrChan{
		errors: make(chan error, b),
	}
}

func (e *ErrChan) Errors() <-chan error {
	return e.errors
}

func (e *ErrChan) Send(err error) {
	select {
	case e.errors <- err:
	default:
	}
}

func (e *ErrChan) Trace(err error) {
	pc, f, l, ok := runtime.Caller(1)
	if !ok {
		e.Send(err)
		return
	}
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = filepath.Base(fn.Name())
	}
	e.Send(errors.Wrapf(err, "%s:%d %s", formatFile(f), l, funcName))
}

func (e *ErrChan) ConsumeErrors(ep ErrorProducer) {
	go func() {
		for err := range ep.Errors() {
			e.Send(err)
		}
	}()
}
