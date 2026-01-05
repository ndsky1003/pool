package pool

import (
	"sync"
	"sync/atomic"
)

type AdaptiveRingPool[T any] struct {
	mu sync.Mutex

	head  int      // 取数位置（高频读写）
	tail  int      // 放数位置（高频读写）
	count int      // 空闲数（高频读写）
	_     [40]byte // 64 - 24 = 40

	hitCount atomic.Int64 // 命中数（高频读，低频写）
	getCount atomic.Int64 // 总获取数（高频写）
	curCap   int          // 当前容量（低频修改）
	_        [40]byte     // 64 - 3*8 = 40

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

// Get 核心：获取对象 + 无锁统计 + 自动学习+伸缩，性能和原生RingBuffer几乎无差别
func (p *AdaptiveRingPool[T]) Get() T {
	// 1. 原子统计：总获取数+1，无锁，零损耗
	p.getCount.Add(1)

	p.mu.Lock()
	defer p.mu.Unlock()

	// 2. 有空闲对象，复用，命中数+1
	if p.count > 0 {
		obj := p.buffer[p.head]
		p.head = (p.head + 1) % p.curCap
		p.count--
		p.hitCount.Add(1)
		return obj
	}

	// 3. 无空闲对象，新建
	return p.New()
}

// Put 核心：放回对象 + 触发自动学习+伸缩逻辑，核心逻辑都在这里
func (p *AdaptiveRingPool[T]) Put(obj T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 1. 队列未满，放回对象
	if p.count < p.curCap {
		p.buffer[p.tail] = obj
		p.tail = (p.tail + 1) % p.curCap
		p.count++
	}
	// 队列已满，直接丢弃，避免内存溢出

	// 2. 核心：自动学习+自适应伸缩，只在Put时触发，频率极低，无性能损耗
	p.autoScale()
}

// autoScale 自动学习+扩容缩容核心逻辑，极简，无复杂计算，锁内执行，耗时可忽略
func (p *AdaptiveRingPool[T]) autoScale() {
	// 总获取数为0，无需伸缩
	total := p.getCount.Load()
	if total == 0 {
		return
	}

	// 计算命中率
	hitRate := float64(p.hitCount.Load()) / float64(total)

	// 情况1：命中率过高 → 忙时，扩容
	if hitRate > p.opt.HitRateHigh && p.curCap < p.opt.MaxCapacity {
		newCap := int(float64(p.curCap) * p.opt.ScaleUpFactor)
		newCap = min(newCap, p.opt.MaxCapacity)
		p.resize(newCap)
		return
	}

	// 情况2：命中率过低 → 闲时，缩容
	if hitRate < p.opt.HitRateLow && p.curCap > p.opt.MinCapacity {
		newCap := int(float64(p.curCap) * p.opt.ScaleDownFactor)
		newCap = max(newCap, p.opt.MinCapacity)
		p.resize(newCap)
		return
	}

	// 情况3：命中率适中，不做任何操作，维持当前容量
}

// resize 环形队列的扩容/缩容实现，最优写法，无内存浪费，性能极致
func (p *AdaptiveRingPool[T]) resize(newCap int) {
	if newCap <= 0 || newCap == p.curCap {
		return
	}

	// 新建新容量的数组
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
	p.curCap = newCap

	// 重置统计，开始新一轮的自动学习
	p.hitCount.Store(0)
	p.getCount.Store(0)
}
