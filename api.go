package pool

import (
	"reflect"
	"sync"
)

type typekey[T any] struct{}

var m sync.Map

func Regist[T any](newFunc func() any) {
	if newFunc == nil {
		panic("pool: newFunc cannot be nil")
	}
	key := typekey[T]{}
	m.Store(key, &sync.Pool{New: newFunc})
}

func Get[T any]() (T, bool) {
	var zero T
	key := typekey[T]{}

	poolVal, ok := m.Load(key)
	if !ok {
		return zero, false
	}

	pool, ok := poolVal.(*sync.Pool)
	if !ok {
		return zero, false
	}

	obj := pool.Get()
	result, ok := obj.(T)
	return result, ok
}

func MustGet[T any]() T {
	key := typekey[T]{}

	poolVal, ok := m.Load(key)
	if !ok {
		panic("pool: type " + getTypeName[T]() + " not registered")
	}

	pool, ok := poolVal.(*sync.Pool)
	if !ok {
		panic("pool: type " + getTypeName[T]() + " registered but not a sync.Pool")
	}

	obj := pool.Get()
	result, ok := obj.(T)
	if !ok {
		panic("pool: object type mismatch for " + getTypeName[T]())
	}
	return result
}

func Put[T any](obj T) {
	key := typekey[T]{}

	poolVal, ok := m.Load(key)
	if !ok {
		return
	}

	pool, ok := poolVal.(*sync.Pool)
	if !ok {
		return
	}
	pool.Put(obj)
}

func Unregist[T any]() {
	key := typekey[T]{}
	m.Delete(key)
}

func getTypeName[T any]() string {
	var zero T
	return string(reflect.TypeOf(zero).String())
}
