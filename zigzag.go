// 整数的压缩算法
package util

func ZigZag32(n int32) int32 {
	return (n >> 1) ^ (n << 31)
}

func ZigZag64(n int64) int64 {
	return (n >> 1) ^ (n << 63)
}

func FromZigZag32(n int32) int32 {
	return int32(uint32(n)>>1) ^ -(n & 1)
}

func FromZigZag64(n int64) int64 {
	return int64(uint64(n)>>1) ^ -(n & 1)
}
