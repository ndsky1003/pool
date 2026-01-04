package pool

type Option struct {
	// 默认最小容量：保底的空闲buffer数，缩容不会低于这个值，内存占用极低
	MinCapacity *int
	// 默认最大容量：扩容不会超过这个值，彻底杜绝内存泄漏，个人项目512足够用
	MaxCapacity *int
	// 扩容因子：忙时扩容1.2倍，平缓扩容，无性能波动
	ScaleUpFactor *float32
	// 缩容因子：闲时缩容0.8倍，平缓缩容，保留足够的空闲buffer
	ScaleDownFactor *float32
	// 命中率阈值：>0.8扩容，<0.2缩容，业界通用最优值，不用调
	HitRateHigh *float32
	HitRateLow  *float32
}

func Options() *Option {
	return &Option{}
}

func (o *Option) SetMinCapacity(delta int) *Option {
	o.MinCapacity = &delta
	return o
}

func (o *Option) SetMaxCapacity(delta int) *Option {
	o.MaxCapacity = &delta
	return o
}

func (o *Option) SetScaleUpFactor(delta float32) *Option {
	o.ScaleUpFactor = &delta
	return o
}

func (o *Option) SetScaleDownFactor(delta float32) *Option {
	o.ScaleDownFactor = &delta
	return o
}

func (o *Option) SetHitRateHigh(delta float32) *Option {
	o.HitRateHigh = &delta
	return o
}

func (o *Option) SetHitRateLow(delta float32) *Option {
	o.HitRateLow = &delta
	return o
}

func (o *Option) merge(delta *Option) {
	if o == nil || delta == nil {
		return
	}
	resolveOption(&o.MinCapacity, delta.MinCapacity)
	resolveOption(&o.MaxCapacity, delta.MaxCapacity)
	resolveOption(&o.ScaleUpFactor, delta.ScaleUpFactor)
	resolveOption(&o.ScaleDownFactor, delta.ScaleDownFactor)
	resolveOption(&o.HitRateHigh, delta.HitRateHigh)
	resolveOption(&o.HitRateLow, delta.HitRateLow)
}

func (o Option) Merge(opts ...*Option) Option {
	for _, opt := range opts {
		o.merge(opt)
	}
	return o
}
