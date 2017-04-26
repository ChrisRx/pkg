package options

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type Mode int

const (
	PanicMode Mode = iota
	PicnicMode
)

var (
	ErrInvalidType = errors.New("invalid type")
)

type Defaults map[string]interface{}
type Opts map[string]interface{}
type Option interface{}

// Opt represents a single key/value option
type Opt struct {
	Key   string
	Value interface{}
}

type Options struct {
	defaults map[string]interface{}
	options  map[string]interface{}

	panic  bool
	picnic bool
}

func New(opts ...Option) *Options {
	if len(opts) == 0 {
		return newOptions()
	}
	switch t := opts[0].(type) {
	case Opts:
		return newFromOpts(t)
	case *Options:
		return t
	case Options:
		return &t
	case Opt:
		return newFromOpt(opts)
	case []Option:
		return New(t[0])
	default:
		//fmt.Printf("Yup1\n")
		panic(fmt.Errorf("unknown option type received, '%v'", reflect.TypeOf(opts[0])))
	}
}

func newOptions() *Options {
	return &Options{
		defaults: make(map[string]interface{}),
		options:  make(map[string]interface{}),
		panic:    true,
	}
}

func newFromOpt(opts []Option) *Options {
	o := newOptions()
	options := make([]Opt, 0)
	for _, o := range opts {
		opt, ok := o.(Opt)
		if !ok {
			panic(fmt.Errorf("unknown option type received, '%v', expected 'options.Option'", reflect.TypeOf(o)))
		}
		options = append(options, opt)
	}
	o.Apply(options...)
	return o
}

func newFromOpts(opts Opts) *Options {
	o := newOptions()
	options := make([]Opt, 0)
	for k, v := range opts {
		options = append(options, Opt{k, v})
	}
	o.Apply(options...)
	return o
}

// GetDefaults returns a new Options from the initial defaults provided
// to the current Options
func (o *Options) GetDefaults() *Options {
	opts := make([]Opt, 0)
	for k, v := range o.defaults {
		opts = append(opts, Opt{k, v})
	}
	return newFromOpts(o.defaults)
}

func (o *Options) checkType(key string, value interface{}) {
	if val, ok := o.options[key]; ok {
		t1 := reflect.TypeOf(val)
		t2 := reflect.TypeOf(value)
		if t1 != t2 {
			// TODO: check if t1 is an interface that implements t2, or maybe not
			panic(fmt.Sprintf("Option '%s' must be value of type '%s', received type '%s'\n", t1, key, t2))
		}
	}
	if val, ok := o.defaults[key]; ok && val != nil {
		t1 := reflect.TypeOf(val)
		t2 := reflect.TypeOf(value)
		if t1 != t2 {
			// TODO: check if t1 is an interface that implements t2, or maybe not
			panic(fmt.Sprintf("Option '%s' must be value of type '%s', received type '%s'\n", t1, key, t2))
		}
	}
}

// TODO: variadic option(s): InvalidateExisting
func (o *Options) SetDefaults(opts Defaults) {
	for k, v := range opts {
		o.defaults[k] = v
		o.checkType(k, v)
	}
}

func (o *Options) WithDefaults(opts Defaults) *Options {
	o.SetDefaults(opts)
	return o
}

func (o *Options) With(modes ...Mode) *Options {
	for _, m := range modes {
		switch m {
		case PanicMode:
			o.panic = true
		case PicnicMode:
			o.panic = false
		}
	}
	return o
}

func (o *Options) Apply(opts ...Opt) {
	for _, opt := range opts {
		if len(o.defaults) > 0 {
			if _, ok := o.defaults[opt.Key]; !ok {
				panic(fmt.Sprintf("Invalid option '%s'", opt.Key))
			}
		}
		o.checkType(opt.Key, opt.Value)
		o.options[opt.Key] = opt.Value
	}
}

func (o *Options) Replace(opts ...Opt) {
	for _, opt := range opts {
		o.options[opt.Key] = opt.Value
	}
}

// get is called internally when the key lookup should not cause a panic
func (o *Options) get(key string) (interface{}, bool) {
	if val, ok := o.options[key]; ok {
		return val, true
	}
	if val, ok := o.defaults[key]; ok {
		return val, true
	}
	return nil, false
}

func (o *Options) Get(key string) (interface{}, bool) {
	if val, ok := o.get(key); ok {
		return val, true
	}
	if o.panic {
		panic(fmt.Errorf("unable to find value for '%s'", key))
	}
	return nil, false
}

func (o *Options) GetObject(key string, v interface{}) bool {
	if val, ok := o.Get(key); ok {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(val))
		return true
	}
	return false
}

// TODO: GetStringMap

func (o *Options) GetBool(key string) (v bool) {
	if val, ok := o.Get(key); ok {
		v, _ = val.(bool)
	}
	return
}

func (o *Options) GetDuration(key string) (v time.Duration) {
	if val, ok := o.Get(key); ok {
		v, _ = val.(time.Duration)
	}
	return
}

func (o *Options) GetInt(key string) (v int) {
	if val, ok := o.Get(key); ok {
		v, _ = val.(int)
	}
	return
}

func (o *Options) GetString(key string) (v string) {
	if val, ok := o.Get(key); ok {
		switch t := val.(type) {
		case string:
			return t
		case fmt.Stringer:
			return t.String()
		default:
			panic(fmt.Errorf("unable to convert type '%v' into string", reflect.TypeOf(val)))
		}
	}
	return
}

func (o *Options) GetTime(key string) (v time.Time) {
	if val, ok := o.Get(key); ok {
		v, _ = val.(time.Time)
	}
	return
}

func (o *Options) GetUrl(key string) (v *url.URL) {
	if val, ok := o.Get(key); ok {
		switch t := val.(type) {
		case *url.URL:
			return t
		default:
			panic(fmt.Errorf("unable to convert type '%v' into *url.URL", reflect.TypeOf(val)))
		}
	}
	return
}

func (o *Options) Set(key string, v interface{}) {
	o.Apply(Opt{key, v})
}

func (o *Options) All() []Opt {
	opts := make([]Opt, 0)
	for k, v := range o.options {
		opts = append(opts, Opt{k, v})
	}
	return opts
}

func (o *Options) String() string {
	t := make(map[string]interface{})
	for k, v := range o.defaults {
		t[k] = v
	}
	for k, v := range o.options {
		t[k] = v
	}
	var output string
	for k, v := range o.Iter() {
		output += fmt.Sprintf("%s: %v\n", k, v)
	}
	return output
}

func (o *Options) Iter() map[string]interface{} {
	t := make(map[string]interface{})
	for k, v := range o.defaults {
		t[k] = v
	}
	for k, v := range o.options {
		t[k] = v
	}
	return t
}
