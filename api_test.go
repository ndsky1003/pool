package pool

import (
	"sync"
	"testing"
)

type TT struct {
	Name string
}

func BenchmarkGet(b *testing.B) {
	Regist[*TT](func() any {
		return &TT{}
	})
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			if v, b := Get[*TT](); b {
				Put(v)
			}
		}
	})

}

func BenchmarkGet2(b *testing.B) {
	pool := sync.Pool{
		New: func() any {
			return &TT{}
		},
	}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v := pool.Get()
			pool.Put(v)
		}
	})

}

func BenchmarkGet3(b *testing.B) {
	pool := New(func() *TT {
		return &TT{}
	})
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v := pool.Get()
			pool.Put(v)
		}
	})

}
