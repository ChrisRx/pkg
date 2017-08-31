package options_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ChrisRx/pkg/options"
)

func TestDefaultOptionPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test did not panic")
		}
	}()
	opts := options.New().WithDefaults(options.Defaults{
		"port": 8080,
	})
	opts.Apply(options.Opt{"port", "eightyeighty"})
}

func TestOptionSetDefaultsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test did not panic")
		}
	}()
	opts := options.New()
	opts.Apply(options.Opt{"port", "eightyeighty"})
	opts.SetDefaults(options.Defaults{
		"port": 8080,
	})
}

func TestInvalidOptionPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test did not panic")
		}
	}()
	opts := options.New().WithDefaults(options.Defaults{
		"port": 8080,
	})
	opts.Apply(options.Opt{"host", "0.0.0.0"})
}

func TestGetObject(t *testing.T) {
	opts := options.New().WithDefaults(options.Defaults{
		"host": "0.0.0.0",
		"port": 8080,
	})
	opts.Apply(options.Opt{"host", "0.0.0.0"})
	var obj string
	opts.GetObject("host", &obj)
}

// TODO: this test needs
func TestGetObjectInterface(t *testing.T) {
	opts := options.New().WithDefaults(options.Defaults{
		"host": "0.0.0.0",
		"port": 8080,
	})
	opts.Apply(options.Opt{"host", "0.0.0.0"})
	var obj string
	opts.GetObject("host", &obj)
}

func TestDefaultNil(t *testing.T) {
	opts := options.New()
	opts.SetDefaults(options.Defaults{
		"port": nil,
	})
	opts.Apply(options.Opt{"port", 8080})
}

func TestOptionSetOverride(t *testing.T) {
	opts := options.New().WithDefaults(options.Defaults{
		"port": "8080",
	})
	opts.Replace(options.Opt{"port", 8080})
}

func TestOptionGetObjectInterface(t *testing.T) {
	opts := options.New().WithDefaults(options.Defaults{
		"timeout": 300 * time.Second,
	})
	var obj fmt.Stringer
	if ok := opts.GetObject("timeout", &obj); !ok {
		t.Errorf("could not get object")
	}
}

func TestOptionPanicMode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test did not panic")
		}
	}()
	opts := options.New().With(options.PanicMode)
	opts.GetString("host")
}

func TestOptionInitalizationOpt(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("test did panic")
		}
	}()
	opts := options.New(
		options.Opt{"host", "0.0.0.0"},
		options.Opt{"port", 8080},
	)
	opts.GetString("host")
}

func NewObject(opts ...options.Option) *options.Options {
	o := options.New(opts)
	return o
}

func TestOptionInitalizationNew(t *testing.T) {
	opts := options.New(
		options.Opt{"host", "0.0.0.0"},
		options.Opt{"port", 8080},
	)
	o := options.New(opts)
	v := reflect.ValueOf(o)
	if _, ok := v.Interface().(*options.Options); !ok {
		t.Errorf("expected type '*options.Options', received %v", reflect.TypeOf(o))
	}
}

func TestOptionInitalizationNewObjectUnwrapSlice(t *testing.T) {
	opts := options.New(
		options.Opt{"host", "0.0.0.0"},
		options.Opt{"port", 8080},
	)
	newObject := func(opts ...options.Option) *options.Options {
		return options.New(opts)
	}
	o := newObject(opts)
	v := reflect.ValueOf(o)
	if _, ok := v.Interface().(*options.Options); !ok {
		t.Errorf("expected type '*options.Options', received %v", reflect.TypeOf(o))
	}
}

func ExampleOptions() {
	opts := options.New().WithDefaults(options.Defaults{
		"host":    "0.0.0.0",
		"port":    8080,
		"timeout": 300 * time.Second,
	})
	fmt.Printf("Serving on '%s:%d'\n", opts.GetString("host"), opts.GetInt("port"))
	// Output: Serving on '0.0.0.0:8080'
}

func ExampleOptionsPrint() {
	opts := options.New(options.Opts{
		"host":    "0.0.0.0",
		"port":    8080,
		"timeout": 300 * time.Second,
	})
	fmt.Println(opts.String())
}
