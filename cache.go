package goBase

import (
	"container/list"
	"context"
	"sync"
)

// GRedis 连接池
//https://github.com/go-redis/redis
type GRedis struct {
	sync.Locker
	MaxOpenCons  int
	MaxIdleConns int
	Name         string
	FreeCons     *list.List
	Cxt          context.Context
}

// GetCon 获取一个连接
func (a *GRedis) GetCon() {
	if a.FreeCons.Len() > 0 {
		a.Lock()
		defer a.Unlock()
		a.FreeCons.Front()
	}
	//return nil
}

// Close 是否一个连接
func (a *GRedis) Close() {}
