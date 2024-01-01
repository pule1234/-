package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func Test_ExpiredLock(t *testing.T) {
	lock := NewExpiredLock()
	lock.Lock(1)
	<-time.After(time.Duration(1) * time.Second)
	lock.Lock(0) // 锁被重复加的话会陷入死锁
	if err := lock.Unlock(); err != nil {
		t.Error(err)
	}
}

type ExpiredLock struct {
	mutex sync.Mutex // 主流程锁

	processMutex sync.Mutex // 辅助锁  保证加锁的原子性

	owner string // 标识锁的归属

	stop context.CancelFunc // 异步goroutine生命周期终止控制器
}

func NewExpiredLock() *ExpiredLock {
	return &ExpiredLock{}
}

func (e *ExpiredLock) Lock(expiredSeconds int) {
	e.mutex.Lock()
	e.processMutex.Lock()
	defer e.processMutex.Unlock()
	token := GetCurrentProcessAndGoroutineIDStr()
	e.owner = token

	if expiredSeconds <= 0 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.stop = cancel
	//保证在达到过期时长之后，执行解锁操作
	go func() {
		select {
		case <-time.After(time.Duration(expiredSeconds) * time.Second):
			e.unLock(token) // 解锁
		case <-ctx.Done():
		}
	}()
}

func (e *ExpiredLock) Unlock() error {
	token := GetCurrentProcessAndGoroutineIDStr()
	return e.unLock(token)
}

// 若在lock方法中的异步协程中调用  Unlock执行GetCurrentProcessAndGoroutineIDStr（）方法的话
// 得到的值和原始的值永远是不一样的
func (e *ExpiredLock) unLock(token string) error {
	e.processMutex.Lock()
	defer e.processMutex.Unlock()

	if token != e.owner {
		return errors.New("not your lock !!!")
	}

	e.owner = ""
	// 停止异步goroutine
	if e.stop != nil {
		e.stop()
	}
	e.mutex.Unlock()
	return nil
}
