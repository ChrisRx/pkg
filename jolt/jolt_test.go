package jolt_test

import (
	"os"
	"testing"
	"time"

	"github.com/ChrisRx/pkg/jolt"
)

func TestInvalidPrintArgsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test did not panic")
		}
	}()
	j := jolt.New(os.Stdout)
	j.Print(0)
}

func ExampleJoltFields() {
	j := jolt.New(os.Stdout)
	j.TimeFunc = func() time.Time {
		t, _ := time.Parse("2006", "2006")
		return t
	}
	j.Print(jolt.Fields{
		"msg": "jolt'n like a sultan",
	})
	//Output: {"msg":"jolt'n like a sultan","ts":"2006-01-01T00:00:00Z"}
}

func ExampleJoltPrintf() {
	j := jolt.New(os.Stdout)
	j.TimeFunc = func() time.Time {
		t, _ := time.Parse("2006", "2006")
		return t
	}
	j.Print("%s'n like a sultan", "jolt")
	//Output: 2006-01-01T00:00:00Z - jolt'n like a sultan
}
