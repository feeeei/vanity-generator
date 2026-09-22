package context

import (
	"sync"
	"sync/atomic"
)

type Context struct {
	power int64
	done  int64
	wait  *sync.WaitGroup
}

func NewContext() *Context {
	wait := &sync.WaitGroup{}
	wait.Add(1)
	return &Context{
		wait: wait,
	}
}

func (c *Context) IsDone() bool {
	return atomic.LoadInt64(&c.done) != 0
}

// Finish 结束等待，多个 worker 同时调用时只有第一次生效
func (c *Context) Finish() {
	if atomic.CompareAndSwapInt64(&c.done, 0, 1) {
		c.wait.Done()
	}
}

func (c *Context) Increment() {
	atomic.AddInt64(&c.power, 1)
}

func (e *Context) Reset() int64 {
	return atomic.SwapInt64(&e.power, 0)
}

func (c *Context) Wait() {
	c.wait.Wait()
}
