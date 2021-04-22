// 生成随机字符串
package util

import (
	"math/rand"
	"strings"
	"time"
)

var (
	randBytes []byte                                        // 随机字符
	random    = rand.New(rand.NewSource(time.Now().Unix())) // 随机数
)

func init() {
	for i := '0'; i <= '9'; i++ {
		randBytes = append(randBytes, byte(i))
	}
	for i := 'a'; i <= 'z'; i++ {
		randBytes = append(randBytes, byte(i))
	}
	for i := 'A'; i <= 'Z'; i++ {
		randBytes = append(randBytes, byte(i))
	}
}

// 设置RandomString产生的字符，不是同步的
func SetRandomBytes(b []byte) {
	randBytes = make([]byte, len(b))
	copy(randBytes, b)
}

// 返回随机字符串"a-z 0-9 A-Z"，n:长度
func RandomString(n int) string {
	var str strings.Builder
	for i := 0; i < n; i++ {
		str.WriteByte(randBytes[random.Intn(len(randBytes))])
	}
	return str.String()
}

// 随机数字字符串"0-9"，n:长度
func RandomNumber(n int) string {
	var str strings.Builder
	for i := 0; i < n; i++ {
		str.WriteByte(randBytes[random.Intn(10)])
	}
	return str.String()
}
