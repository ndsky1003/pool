package pool

import (
	"sync"
)

type AdaptiveRingPool[T any] struct {
	mu sync.Mutex

	head     int // 取数位置（高频读写）
	tail     int // 放数位置（高频读写）
	count    int // 空闲数（高频读写）
	putcount int
	_        [32]byte // 64 - 24 = 40

	curCap int      // 当前容量（低频修改）
	_      [56]byte // 64 - 8 = 56

	opt Option
	_   [16]byte

	New    func() T // 创建函数（初始化后不变）
	buffer []T      // 环形队列（低频大尺寸访问）
	_      [32]byte // 64 - 8 - 24  = 32
}

// NewAdaptiveRingPool 创建自适应环形池，个人项目无脑用这个，默认配置足够
func NewAdaptiveRingPool[T any](newFunc func() T, opts ...OptionFunc) *AdaptiveRingPool[T] {
	opt := DefaultOptions()
	for _, f := range opts {
		f(&opt)
	}

	if opt.MinCapacity <= 0 {
		opt.MinCapacity = 32
	}

	if opt.MaxCapacity < opt.MinCapacity {
		opt.MaxCapacity = opt.MinCapacity
	}

	o := &AdaptiveRingPool[T]{
		buffer: make([]T, opt.MinCapacity),
		curCap: opt.MinCapacity,
		opt:    opt,
		New:    newFunc,
	}
	return o
}

// ------------------------ 你的原始 Get 方法，完全不变 ------------------------
func (p *AdaptiveRingPool[T]) Get() T {
	p.mu.Lock()
	// 2. 有空闲对象，复用，命中数+1
	if p.count > 0 {
		obj := p.buffer[p.head]
		// 核心优化：用位运算替代取模 (head+1) % curCap → 仅当curCap是2的幂次时可用
		p.head = (p.head + 1) & (p.curCap - 1)
		p.count--
		p.mu.Unlock()
		return obj
	}
	p.mu.Unlock()

	// 3. 无空闲对象，新建
	return p.New()
}

// ------------------------ 你的原始 Put 方法，完全不变 ------------------------
func (p *AdaptiveRingPool[T]) Put(obj T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.putcount++

	if p.count < p.curCap {
		p.buffer[p.tail] = obj
		// 核心优化：用位运算替代取模 (tail+1) % curCap → 仅当curCap是2的幂次时可用
		p.tail = (p.tail + 1) & (p.curCap - 1)
		p.count++
	}

	if p.putcount&0xf == 0xf {
		p.autoScale()
	}
}

func (p *AdaptiveRingPool[T]) autoScale() {
	usageRatio := calcPercentPowerOfTwo(int32(p.count), log2(p.curCap))
	usageRatio = 100 - usageRatio
	if usageRatio >= p.opt.ScaleUpRatio && p.curCap < p.opt.MaxCapacity {
		newCap := p.curCap << 1
		newCap = min(newCap, p.opt.MaxCapacity)
		p.resize(newCap)
		return
	}

	if usageRatio <= p.opt.ScaleDownRatio && p.curCap > p.opt.MinCapacity {
		newCap := p.curCap >> 1
		newCap = max(newCap, p.opt.MinCapacity)
		p.resize(newCap)
		return
	}

}

func (p *AdaptiveRingPool[T]) resize(newCap int) {
	if newCap <= 0 || newCap == p.curCap {
		return
	}

	// 新建新容量的数组 → 新容量是2的幂次
	newBuf := make([]T, newCap)
	// 计算原有效数据的起止索引
	if p.count > 0 {
		end := p.head + p.count
		if end <= p.curCap {
			// 数据未跨环，直接拷贝
			copy(newBuf, p.buffer[p.head:end])
		} else {
			// 数据跨环，分两段拷贝
			copyLen1 := p.curCap - p.head
			copy(newBuf, p.buffer[p.head:])
			copy(newBuf[copyLen1:], p.buffer[:end-p.curCap])
		}
	}

	// 更新队列状态，完成伸缩
	p.buffer = newBuf
	p.head = 0
	p.tail = p.count
	p.putcount = 0
	p.curCap = newCap // newCap 是2的幂次，保证后续位运算可用

}
