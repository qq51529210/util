package util

import (
	"encoding/binary"
	"io"
	"sync"
)

var (
	readBufferPool   sync.Pool
	encodeBufferPool sync.Pool
)

func init() {
	encodeBufferPool.New = func() interface{} {
		return &encodeBuffer{}
	}
	readBufferPool.New = func() interface{} {
		return &readBuffer{}
	}
}

// Use for read a segmented message.
type readBuffer struct {
	b []byte
	n int
}

// Make sure b can stores n bytes data.
func (b *readBuffer) grow(n int) {
	m := b.n + n
	if m > len(b.b) {
		nb := make([]byte, m)
		copy(nb, b.b)
		b.b = nb
	}
}

func (b *readBuffer) ReadN(r io.Reader, n int) error {
	b.grow(n)
	n, err := io.ReadFull(r, b.b[b.n:b.n+n])
	if err != nil {
		return err
	}
	b.n += n
	return nil
}

// Use for encode message.
type encodeBuffer struct {
	b []byte
	n int
}

func (b *encodeBuffer) Reset() {
	b.n = 0
}

// Make sure b can stores n bytes data.
func (b *encodeBuffer) grow(n int) {
	m := b.n + n
	if m > len(b.b) {
		nb := make([]byte, m)
		copy(nb, b.b)
		b.b = nb
	}
}

func (b *encodeBuffer) Put8(n byte) {
	b.grow(1)
	b.b[b.n] = n
	b.n++
}

// Append BigEndian n
func (b *encodeBuffer) Put16(n uint16) {
	b.grow(2)
	binary.BigEndian.PutUint16(b.b[b.n:], n)
	b.n += 2
}

// Append BigEndian n
func (b *encodeBuffer) Put32(n uint32) {
	b.grow(4)
	binary.BigEndian.PutUint32(b.b[b.n:], n)
	b.n += 4
}

// Append BigEndian n
func (b *encodeBuffer) Put64(n uint64) {
	b.grow(8)
	binary.BigEndian.PutUint64(b.b[b.n:], n)
	b.n += 8
}

// Append n random bytes.
func (b *encodeBuffer) PutRandom(n int) {
	b.grow(n)
	random.Read(b.b[b.n:])
	b.n += n
}

// Append bytes.
func (b *encodeBuffer) PutBytes(d []byte) {
	b.grow(len(d))
	b.n += copy(b.b[b.n:], d)
}

// Append string.
func (b *encodeBuffer) PutString(s string) {
	b.grow(len(s))
	b.n += copy(b.b[b.n:], s)
}
