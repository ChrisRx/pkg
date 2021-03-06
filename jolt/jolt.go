package jolt

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/ChrisRx/pkg/errchan"
)

type Fields map[string]interface{}

type FieldFunc func() string

type Jolt struct {
	defaults Fields

	mu sync.Mutex
	w  io.Writer
}

func New(w io.Writer) *Jolt {
	return &Jolt{
		w:        w,
		defaults: make(Fields),
	}
}

func (j *Jolt) With(m Fields) {
	j.defaults = m
}

func (j *Jolt) PrintAll(name string, ep errchan.ErrorProducer) {
	for err := range ep.Errors() {
		j.Print(Fields{name: err.Error()})
	}
}

func (j *Jolt) Print(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	for i, a := range args {
		switch t := a.(type) {
		case Fields:
			j.printFields(t)
			args = append(args[:i], args[i+1:]...)
		}
	}
	if len(args) > 0 {
		var format string
		switch t := args[0].(type) {
		case string:
			format = t
		case fmt.Stringer:
			format = t.String()
		case error:
			format = t.Error()
		default:
			panic(fmt.Errorf("received invalid type '%v' in arguments", reflect.TypeOf(args[0])))
		}
		j.printf(format, args[1:]...)
	}
}

func (j *Jolt) printf(format string, a ...interface{}) {
	j.printFields(Fields{
		"msg": fmt.Sprintf(format, a...),
	})
}

func (j *Jolt) printFields(m Fields) {
	tmp := make(Fields)
	for k, v := range j.defaults {
		switch t := v.(type) {
		case func() string:
			tmp[k] = t()
		case FieldFunc:
			tmp[k] = t()
		default:
			tmp[k] = v
		}
	}
	for k, v := range m {
		tmp[k] = v
	}
	b, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	b = append(b, '\n')
	j.mu.Lock()
	j.w.Write(b)
	j.mu.Unlock()
}
