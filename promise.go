package promise

import (
	"fmt"
	"sync"
)

type function func(interface{}) interface{}
type consumer func(interface{})
type executor func(consumer, consumer)

type status int

const sPENDING status = 0
const sRESOLVED status = 1
const sREJECTED status = 2

type Promise struct {
	stat       status
	value      interface{}
	err        error
	lock       sync.Mutex
	resHandler consumer
	rejHandler consumer
	errHandler consumer
	next       *Promise
}

func New(exec executor) *Promise {
	p := &Promise{}
	exec(func(v interface{}) {
		p.handleRes(v)
	}, func(v interface{}) {
		p.handleRej(v)
	})
	return p
}

func Resolve(v interface{}) *Promise {
	return New(func(res, rej consumer) {
		res(v)
	})
}

func Reject(v interface{}) *Promise {
	return New(func(res, rej consumer) {
		rej(v)
	})
}

func (p *Promise) handleRes(v interface{}) {
	if p.stat != sPENDING {
		return
	}
	p.stat = sRESOLVED
	p.value = v
	if p.next != nil && p.next.resHandler != nil {
		p.next.resHandler(v)
	}
}

func (p *Promise) handleRej(v interface{}) {
	if p.stat != sPENDING {
		return
	}
	p.stat = sREJECTED
	p.value = v
	if p.next != nil && p.next.rejHandler != nil {
		p.next.rejHandler(v)
	}
}

func (p *Promise) handleCatch(err error) {
	if p.stat != sPENDING {
		return
	}
	p.stat = sREJECTED
	p.value = err
	nextP := p.next
	for {
		if nextP != nil {
			errHandler := p.next.errHandler
			if errHandler != nil {
				errHandler(err)
				break
			}
			nextP = nextP.next
		}
	}

}

func (p *Promise) Then(onRes function, onRej function) *Promise {
	newP := &Promise{}
	if p.stat == sPENDING {
		p.next = newP
		newP.resHandler = func(v interface{}) {
			fmt.Println("resh")
			if onRes != nil {
				handleValue(newP, p.value, onRes)
				//newP.handleRes(onRes(v))
			}
		}
		newP.rejHandler = func(v interface{}) {
			fmt.Println("rejh")
			if onRej != nil {
				handleValue(newP, p.value, onRej)
				//newP.handleRes(onRej(v))
			}
		}
	} else if p.stat == sRESOLVED {
		fmt.Println("res")
		if onRes != nil {
			handleValue(newP, p.value, onRes)
			//newP.handleRes(onRes(p.value))
		}
	} else if p.stat == sREJECTED {
		fmt.Println("rej")
		if onRej != nil {
			handleValue(newP, p.value, onRej)
			//newP.handleRes(onRej(p.value))
		}
	}
	return newP
}

func handleValue(p *Promise, input interface{}, fn function) {
	value := fn(input)
	if pp, ok := value.(*Promise); ok {
		pp.Then(func(v interface{}) interface{} {
			p.handleRes(v)
			return v
		}, func(v interface{}) interface{} {
			p.handleRej(v)
			return v
		})
	} else {
		p.handleRes(value)
	}
}

func (p *Promise) Catch(onErr function) *Promise {
	newP := &Promise{}
	p.next = newP
	if p.stat == sPENDING {
		newP.errHandler = func(err interface{}) {
			onErr(err)
		}
	}
	return newP
}
