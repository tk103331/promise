package promise

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type resolver func(interface{}) interface{}
type rejecter func(interface{}) interface{}
type executor func(resolver, rejecter)

type status int

var initial status = 0
var resolved status = 1
var rejected status = 2

type Promise struct {
	stat    status
	data    interface{}
	lock    sync.Mutex
	resFunc resolver
	rejFunc rejecter
}

func (p *Promise) resolve(v interface{}) interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data = v
	p.stat = resolved
	fmt.Println("resolve")
	if p.resFunc != nil {
		p.resFunc(p.data)
	}
	return nil
}

func (p *Promise) reject(v interface{}) interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data = v
	p.stat = rejected
	fmt.Println("reject")
	if p.rejFunc != nil {
		p.rejFunc(p.data)
	}
	return nil
}

func (p *Promise) Then(resolveFunc func(interface{}) interface{}, rejectFunc func(interface{}) interface{}) *Promise {
	p.lock.Lock()
	defer p.lock.Unlock()
	fmt.Println("then")
	var data interface{} = nil
	if p.stat == resolved {
		data = resolveFunc(p.data)
		if pp, ok := data.(*Promise); ok {
			return pp
		} else {
			return Resolve(data)
		}
	} else if p.stat == rejected {
		data = rejectFunc(p.data)
		if pp, ok := data.(*Promise); ok {
			return pp
		} else {
			return Resolve(data)
		}
	} else {
		p.resFunc = resolveFunc
		p.rejFunc = rejectFunc
		return p
	}

}

func Resolve(v interface{}) *Promise {
	return New(func(res resolver, rej rejecter) {
		res(v)
	})
}

func New(exec executor) *Promise {
	p := &Promise{}
	exec(p.resolve, p.reject)
	return p
}
