package util

// var (
// 	bitWriterOr = []uint8{
// 		0b00000000,
// 		0b00000001,
// 		0b00000011,
// 		0b00000111,
// 		0b00001111,
// 		0b00011111,
// 		0b00111111,
// 		0b01111111,
// 		0b11111111,
// 	}
// )

// BitWriter 用于顺序写入 bit
type BitWriter struct {
	// 字节缓存
	b []byte
	// 当前写入的字节
	c byte
	// 当前 bit 的个数
	n int
}

// Reset 重置
func (bw *BitWriter) Reset() {
	bw.b = bw.b[:0]
	bw.n = 0
}

// Bytes 返回当前的缓存, 如果没有写完，补全最后一个
func (bw *BitWriter) Bytes() []byte {
	if bw.n > 0 {
		bw.b = append(bw.b, bw.c)
		bw.n = 0
		bw.c = 0
	}
	return bw.b
}

// Raw 返回当前的缓存
func (bw *BitWriter) Raw() []byte {
	return bw.b
}

// Raw 返回当前的缓存
func (bw *BitWriter) IsCompleted() bool {
	return bw.n == 0
}

// Write8 写入
func (bw *BitWriter) Write8(v uint8, n int) {
	// 判断
	if n > 8 || n < 1 {
		panic("n > 8 or n < 1")
	}
	bw.c |= (v << (8 - n)) >> bw.n
	bw.n += n
	if bw.n > 8 {
		bw.b = append(bw.b, bw.c)
		bw.n -= 8
		bw.c = v << (8 - bw.n)
	} else if bw.n == 8 {
		bw.b = append(bw.b, bw.c)
		bw.n = 0
		bw.c = 0
	}
}

// Write16 写入
func (bw *BitWriter) Write16(v uint16, n int) {

}

// Write32 写入
func (bw *BitWriter) Write32(v uint32, n int) {

}

// Write64 写入
func (bw *BitWriter) Write64(v uint64, n int) {

}
