package promise

import (
	"fmt"
	"testing"
)

func asyncTask(task func()) {
	go task()
}

func TestThenResolve(t *testing.T) {

	New(func(resolve handler, reject handler) {
		asyncTask(func() {
			resolve("resolved value")
		})
	}).Then(func(v interface{}) interface{} {
		fmt.Println("TestThenResolve result1:", v)
		return "another resolved value"
	}, nil).Then(func(v interface{}) interface{} {
		fmt.Println("TestThenResolve result2:", v)
		return nil
	}, nil)
	fmt.Println("...")
}

func TestThenReject(t *testing.T) {
	// New(func(resolve handler, reject handler) {
	// 	asyncTask(func() {
	// 		resolve("rejected value")
	// 	})
	// }).Then(nil, func(v interface{}) interface{} {
	// 	fmt.Println("TestThenReject result:", v)
	// 	return nil
	// })
}
