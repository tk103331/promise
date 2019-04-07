package promise

import (
	"fmt"
	"sync"
)

type handler func(interface{}) interface{}
type executor func(handler, handler)

type status int

var sRESOLVED status = 1
var sREJECTED status = 2

type Promise struct {
	stat status
	data interface{}
	lock sync.Mutex
	res  handler
	rej  handler
	next *Promise
}

func (p *Promise) resolve(v interface{}) interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data = v
	p.stat = sRESOLVED

	if p.res != nil {
		fmt.Printf("%p resolve1\n", p)
		d := p.res(p.data)
		fmt.Printf("%p resolve11\n", p)
		if p.next != nil {
			p.next.resolve(d)
		}
	} else if p.next != nil {
		fmt.Printf("%p resolve2\n", p)
		p.next.resolve(p.data)
		fmt.Printf("%p resolve21\n", p)
	}
	return nil
}

func (p *Promise) reject(v interface{}) interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data = v
	p.stat = sREJECTED

	if p.rej != nil {
		d := p.rej(v)
		if p.next != nil {
			p.next.resolve(d)
		}
	} else if p.next != nil {
		return p.next.reject(p.data)
	}
	return nil
}

func (p *Promise) handleResult(h handler) *Promise {
	result := h(p.data)
	fmt.Printf("%p ", p)
	if pp, ok := result.(*Promise); ok {
		fmt.Printf("%p handleResult1\n", p)
		return pp
	} else {
		fmt.Printf("%p handleResult2\n", p)
		return Resolve(result)
	}
}

func (p *Promise) Then(onResolve handler, onReject handler) *Promise {
	p.lock.Lock()
	defer p.lock.Unlock()
	fmt.Printf("%p ", p)
	if p.stat == sRESOLVED {
		fmt.Println("then res ")
		return p.handleResult(onResolve)
	} else if p.stat == sREJECTED {
		fmt.Println("then rej ")
		return p.handleResult(onReject)
	} else {
		fmt.Println("then else ")
		newP := &Promise{res: onResolve, rej: onReject}
		p.next = newP
		return newP
	}
}

func Resolve(v interface{}) *Promise {
	return New(func(res handler, rej handler) {
		res(v)
	})
}

func New(exec executor) *Promise {
	p := &Promise{}
	exec(p.resolve, p.reject)
	return p
}
