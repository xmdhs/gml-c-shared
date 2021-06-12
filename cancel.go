package main

import (
	"context"
	"sync"
	"sync/atomic"
)

var (
	gcancel = cancel{}
)

type cancel struct {
	sync.Map
	i int64
}

func (c *cancel) add(id int64, f func()) {
	c.Store(id, f)
}

func (c *cancel) del(id int64) {
	c.Delete(id)
}

func (c *cancel) do(id int64) {
	f, ok := c.Load(id)
	if ok {
		ff := f.(func())
		ff()
		c.del(id)
	}
}

func (c *cancel) new() (context.Context, int64) {
	cxt := context.Background()
	cxt, f := context.WithCancel(cxt)
	i := c.getId()
	c.add(i, f)
	return cxt, i
}

func (c *cancel) getId() int64 {
	i := atomic.AddInt64(&c.i, 1)
	return i
}
