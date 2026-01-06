package pool

import "math/bits"

func nextPowerOfTwo(n int32) int32 {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

func calcPercentPowerOfTwo(num int32, totalBits int32) int32 {
	if num < 0 || totalBits <= 0 {
		return 0
	}
	total := int32(1) << totalBits // 分母=2^totalBits
	if num >= total {
		return 100
	}

	// 全程无除法：(num * scale / total) *100 >> scaleBits
	// → 等价于 (num * scale >> totalBits) *100 >> scaleBits
	scaledResult := int64(num) * scale >> totalBits
	percent := (scaledResult * 100) >> scaleBits

	return int32(percent)
}

const scaleBits = 20
const scale = 1 << scaleBits // 1048576

// 预定义log₂查表（覆盖2⁰到2³0）
var log_2_table = [31]int32{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
	20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
}

// Log2Table 查表法计算log₂(n)，仅支持2的幂次
func log2(n int) int32 {
	if n <= 0 || (n&(n-1)) != 0 {
		panic("n must be power of two")
	}
	index := bits.Len(uint(n)) - 1
	return log_2_table[index]
}
