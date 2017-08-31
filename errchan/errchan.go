package errchan

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

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

func (e *ErrChan) Send(errs ...interface{}) error {
	if len(errs) == 0 {
		return nil
	}
	for _, i := range errs {
		switch t := i.(type) {
		case ErrorProducer:
			e.consume(t)
			return nil
		case error:
			select {
			case e.errors <- t:
			default:
			}
			return t
		default:
			panic(fmt.Errorf("unexpected type '%v'", reflect.TypeOf(i)))
		}
	}
	return nil
}

func (e *ErrChan) Trace(err error, args ...interface{}) error {
	pc, f, l, ok := runtime.Caller(+1)
	if !ok {
		return e.Send(fmt.Errorf("Failed Trace: %v", err))
	}
	fp := filepath.Join(filepath.Base(filepath.Dir(f)), filepath.Base(f))
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
	}
	return e.Send(errors.Wrapf(err, "%s:%d %s", fp, l, funcName))
}

func (e *ErrChan) consume(ep ErrorProducer) {
	go func() {
		for err := range ep.Errors() {
			e.Send(err)
		}
	}()
}
