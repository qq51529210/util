package util

import "io"

type BitReader struct {
	w    io.Reader
	buf  byte
	left uint // buf 中还有多少未写入的位
}

// 创建一个 BitReader
func NewBitReader(w io.Reader) *BitReader {
	return &BitReader{w: w}
}
