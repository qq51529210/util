// 因为没有泛型，整数的比较
package util

import "time"

func MaxInt(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxUInt(n1, n2 uint) uint {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxInt16(n1, n2 int16) int16 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxUInt16(n1, n2 uint16) uint16 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxInt32(n1, n2 int32) int32 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxUInt32(n1, n2 uint32) uint32 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxInt64(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxUInt64(n1, n2 uint64) uint64 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func MaxDuration(n1, n2 time.Duration) time.Duration {
	if n1 > n2 {
		return n1
	}
	return n2
}
