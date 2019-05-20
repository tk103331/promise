package promise

import (
	"fmt"
	"strconv"
	"testing"
)

func asyncTask(task func()) {
	go task()
}

func somePromise(n int) []*Promise {
	ps := make([]*Promise, n)
	for i := 0; i < n; i++ {
		func(x int) {
			ps[x] = New(func(resolve consumer, reject consumer) {
				asyncTask(func() {
					resolve("resolved value:" + strconv.Itoa(x))
				})
			})
		}(i)
	}
	return ps
}

func TestResolve(t *testing.T) {
	New(func(resolve consumer, reject consumer) {
		asyncTask(func() {
			resolve("resolved value")
		})
	}).Then(func(v interface{}) interface{} {
		fmt.Println("TestResolve result1:", v)
		return "another resolved value"
	}, nil).Then(func(v interface{}) interface{} {
		fmt.Println("TestThenResolve result2:", v)
		return nil
	}, nil)
	fmt.Println("TestResolve")
}

func TestReject(t *testing.T) {
	New(func(resolve consumer, reject consumer) {
		asyncTask(func() {
			reject("rejected value")
		})
	}).Then(nil, func(v interface{}) interface{} {
		fmt.Println("TestReject result1:", v)
		return "resolve value"
	}).Then(func(v interface{}) interface{} {
		fmt.Println("TestReject result2:", v)
		return nil
	}, nil)
	fmt.Println("TestReject")
}

func TestPromiseValue(t *testing.T) {
	New(func(resolve consumer, reject consumer) {
		asyncTask(func() {
			reject("rejected value")
		})
	}).Then(nil, func(v interface{}) interface{} {
		fmt.Println("TestPromiseValue result1:", v)
		return Resolve("resolve value")
	}).Then(func(v interface{}) interface{} {
		fmt.Println("TestPromiseValue result2:", v)
		return Reject("rejected value")
	}, nil).Then(nil, func(v interface{}) interface{} {
		fmt.Println("TestPromiseValue result3:", v)
		return nil
	})
	fmt.Println("TestPromiseValue")
}

func TestWrap(t *testing.T) {
	warpFunc := Wrap(func() interface{} {
		return "resolve value"
	})

	warpFunc().Then(func(v interface{}) interface{} {
		fmt.Println("TestWrap result:", v)
		return nil
	}, nil)
}

func TestAll(t *testing.T) {
	All(somePromise(5)...).Then(func(v interface{}) interface{} {
		fmt.Println("TestAll result:", v)
		if vs, ok := v.([]interface{}); ok {
			for i := 0; i < len(vs); i++ {
				fmt.Println("\t", vs[i])
			}
		}
		return nil
	}, nil)
}

func TestRace(t *testing.T) {
	Race(somePromise(5)...).Then(func(v interface{}) interface{} {
		fmt.Println("TestRace result:", v)
		return nil
	}, nil)
}
