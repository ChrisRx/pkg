package errchan_test

import (
	"fmt"

	"github.com/ChrisRx/pkg/errchan"
)

func ExampleErrChan() {
	e := errchan.New(100)
	e.Trace(fmt.Errorf("testing"))
	err := <-e.Errors()
	fmt.Println(err)
	//Output: errchan/errchan_test.go:11 errchan_test.ExampleErrChan: testing
}
