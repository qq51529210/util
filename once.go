package util

import (
	"sync"
	"sync/atomic"
)

// Once 保证启动并发状态下，
// 只启动一个协程来处理函数
type Once[T any] struct {
	// 用于等待协程结束
	w sync.WaitGroup
	// 状态，0/1
	s int32
	// 数据
	D T
	// 回调函数
	F func(m T)
}

// Do 尝试启动协程并回调
func (o *Once[T]) Do() {
	if atomic.CompareAndSwapInt32(&o.s, 0, 1) {
		o.w.Add(1)
		go o.routine()
	}
}

// routine 协程中回调函数
func (o *Once[T]) routine() {
	defer func() {
		atomic.StoreInt32(&o.s, 0)
		o.w.Done()
	}()
	// 回调
	o.F(o.D)
}

// Wait 等待结束
func (o *Once[T]) Wait() {
	o.w.Wait()
}
