package pool

import "sync"

type typekey[T any] struct{}

var m sync.Map

func Regist[T any](newFunc func() *T, opts ...OptionFunc) {
	key := typekey[T]{}
	m.Store(key, NewAdaptiveRingPool(newFunc, opts...))
}

func Get[T any]() *T {
	key := typekey[T]{}
	if v, ok := m.Load(key); ok {
		return v.(*AdaptiveRingPool[*T]).Get()
	}
	var zero T
	return &zero
}

func Unregist[T any]() {
	key := typekey[T]{}
	m.Delete(key)
}
