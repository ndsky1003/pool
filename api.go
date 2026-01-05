package pool

import "sync"

type typekey[T any] struct{}

var m sync.Map

func Regist[T any](newFunc func() T, opts ...OptionFunc) {
	key := typekey[T]{}
	m.Store(key, NewAdaptiveRingPool(newFunc, opts...))
}

func Get[T any]() (T, bool) {
	key := typekey[T]{}
	if v, ok := m.Load(key); ok {
		if vv, ok1 := v.(*AdaptiveRingPool[T]); ok1 {
			return vv.Get(), true
		}
	}
	var zero T
	return zero, false
}

func MustGet[T any]() T {
	key := typekey[T]{}
	if v, ok := m.Load(key); ok {
		if vv, ok1 := v.(*AdaptiveRingPool[T]); ok1 {
			return vv.Get()
		}
	}
	panic("pool: type not registered")
}

func Put[T any](obj T) {
	key := typekey[T]{}
	if v, ok := m.Load(key); ok {
		if vv, ok1 := v.(*AdaptiveRingPool[T]); ok1 {
			vv.Put(obj)
		}
	}
}

func Unregist[T any]() {
	key := typekey[T]{}
	m.Delete(key)
}
