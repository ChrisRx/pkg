package jolt

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"
)

type Fields map[string]interface{}

type Jolt struct {
	TimeFunc func() time.Time

	mu sync.Mutex
	w  io.Writer
}

func New(w io.Writer) *Jolt {
	return &Jolt{
		w: w,
		TimeFunc: func() time.Time {
			return time.Now().UTC()
		},
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
		format, ok := args[0].(string)
		if !ok {
			panic(fmt.Errorf("received invalid type '%v' in arguments", reflect.TypeOf(args[0])))
		}
		j.printf(format, args[1:]...)
	}
}

func (j *Jolt) printf(format string, a ...interface{}) {
	ts := j.TimeFunc().Format(time.RFC3339Nano)
	format = fmt.Sprintf("%s - %s", ts, format)
	j.mu.Lock()
	fmt.Fprintf(j.w, format, a...)
	j.mu.Unlock()
}

func (j *Jolt) printFields(m Fields) {
	if _, ok := m["ts"]; !ok {
		m["ts"] = j.TimeFunc().Format(time.RFC3339Nano)
	}
	b, _ := json.Marshal(m)
	b = append(b, '\n')
	j.mu.Lock()
	j.w.Write(b)
	j.mu.Unlock()
}
