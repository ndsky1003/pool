package pool

import (
	"fmt"
	"testing"
)

func BenchmarkCalcPercentPowerOfTwo(b *testing.B) {
	tests := []struct {
		num       int32
		totalBits int32
	}{
		{num: 0, totalBits: 10},
		{num: 512, totalBits: 10},
		{num: 1024, totalBits: 10},
		{num: 2048, totalBits: 11},
		{num: 4096, totalBits: 12},
		{num: 8192, totalBits: 13},
		{num: 16384, totalBits: 14},
		{num: 32768, totalBits: 15},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("num=%d_totalBits=%d", tt.num, tt.totalBits), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				calcPercentPowerOfTwo(tt.num, tt.totalBits)
			}
		})
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		input    int32
		expected int32
	}{
		{input: -5, expected: 1},
		{input: 0, expected: 1},
		{input: 1, expected: 1},
		{input: 2, expected: 2},
		{input: 3, expected: 4},
		{input: 5, expected: 8},
		{input: 16, expected: 16},
		{input: 17, expected: 32},
		{input: 31, expected: 32},
		{input: 64, expected: 64},
	}

	for _, tt := range tests {
		result := nextPowerOfTwo(tt.input)
		if result != tt.expected {
			t.Errorf("nextPowerOfTwo(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func BenchmarkNextPowerOfTwo(b *testing.B) {
	tests := []int32{-5, 0, 1, 2, 3, 5, 16, 17, 31, 64, 100, 255, 512, 1023, 2048}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("input=%d", tt), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				nextPowerOfTwo(tt)
			}
		})
	}
}

func BenchmarkLog2(b *testing.B) {
	tests := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("input=%d", tt), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				log2(tt)
			}
		})
	}
}

func TestCalcPercentPowerOfTwo(t *testing.T) {
	tests := []struct {
		num       int32
		totalBits int32
		expected  int32
	}{
		{num: -1, totalBits: 10, expected: 0},
		{num: 0, totalBits: 10, expected: 0},
		{num: 512, totalBits: 10, expected: 50},
		{num: 1024, totalBits: 10, expected: 100},
		{num: 2048, totalBits: 11, expected: 100},
		{num: 256, totalBits: 10, expected: 25},
		{num: 768, totalBits: 10, expected: 75},
	}

	for _, tt := range tests {
		result := calcPercentPowerOfTwo(tt.num, tt.totalBits)
		if result != tt.expected {
			t.Errorf("calcPercentPowerOfTwo(%d, %d) = %d; want %d", tt.num, tt.totalBits, result, tt.expected)
		}
	}
}

func TestLog2(t *testing.T) {
	tests := []struct {
		input    int
		expected int32
	}{
		{input: 1, expected: 0},
		{input: 2, expected: 1},
		{input: 4, expected: 2},
		{input: 8, expected: 3},
		{input: 16, expected: 4},
		{input: 32, expected: 5},
		{input: 64, expected: 6},
		{input: 128, expected: 7},
		{input: 256, expected: 8},
	}

	for _, tt := range tests {
		result := log2(tt.input)
		if result != tt.expected {
			t.Errorf("log2(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}

	// Test panic for non-power of two
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("log2 did not panic for non-power of two input")
		}
	}()
	log2(3) // This should panic
}
