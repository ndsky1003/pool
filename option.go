package pool

type Option struct {
	// 默认最小容量：保底的空闲buffer数，缩容不会低于这个值，内存占用极低
	MinCapacity int
	// 默认最大容量：扩容不会超过这个值，彻底杜绝内存泄漏，个人项目512足够用
	MaxCapacity int
	// 扩容因子：忙时扩容1.2倍，平缓扩容，无性能波动
	ScaleUpRatio int32
	// 缩容因子：闲时缩容0.8倍，平缓缩容，保留足够的空闲buffer
	ScaleDownRatio int32
}

// 第一步：定义选项函数类型（核心）
type OptionFunc func(*Option)

// 第三步：返回默认配置（值类型，上游决定是否逃逸）
func DefaultOptions() Option {
	return Option{
		MinCapacity:    32,
		MaxCapacity:    512,
		ScaleUpRatio:   80,
		ScaleDownRatio: 20,
	}
}

// 第四步：定义选项函数（显式配置才调用，天然区分“是否设置”）
func WithMinCapacity(v int) OptionFunc {
	return func(o *Option) {
		if v > 0 { // 防御性校验
			o.MinCapacity = int(nextPowerOfTwo(int32(32)))
		}
	}
}

func WithMaxCapacity(v int) OptionFunc {
	return func(o *Option) {
		if v > 0 {
			o.MaxCapacity = int(nextPowerOfTwo(int32(32)))
		}
	}
}

func WithScaleUpRatio(v int32) OptionFunc {
	return func(o *Option) {
		o.ScaleUpRatio = v
	}
}

func WithScaleDownRatio(v int32) OptionFunc {
	return func(o *Option) {
		o.ScaleDownRatio = v
	}
}
