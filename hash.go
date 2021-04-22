// 哈希缓存持
package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"sync"
)

var (
	md5Pool  sync.Pool
	sha1Pool sync.Pool
)

func init() {
	md5Pool.New = func() interface{} {
		return &hashBuffer{
			hash: md5.New(),
			buf:  make([]byte, 0, md5.Size*2),
			sum:  make([]byte, md5.Size),
		}
	}
	sha1Pool.New = func() interface{} {
		return &hashBuffer{
			hash: sha1.New(),
			buf:  make([]byte, 0, sha1.Size*2),
			sum:  make([]byte, sha1.Size),
		}
	}
}

// 对字符串s作md5的运算，然后返回16进制的字符串值
func MD5(s string) string {
	h := md5Pool.Get().(*hashBuffer)
	s = h.Hash(s)
	md5Pool.Put(h)
	return s
}

// 对字符串s作sha1的运算，然后返回16进制的字符串值
func SHA1(s string) string {
	h := sha1Pool.Get().(*hashBuffer)
	s = h.Hash(s)
	sha1Pool.Put(h)
	return s
}

// 用于做hash运算的缓存
type hashBuffer struct {
	hash hash.Hash
	buf  []byte
	sum  []byte
}

// hash并转成15进制
func (h *hashBuffer) Hash(s string) string {
	h.buf = h.buf[:0]
	h.buf = append(h.buf, s...)
	h.hash.Reset()
	h.hash.Write(h.buf)
	h.hash.Sum(h.sum[:0])
	h.buf = h.buf[:h.hash.Size()*2]
	hex.Encode(h.buf, h.sum)
	return string(h.buf)
}
