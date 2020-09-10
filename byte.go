package common

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	KBits = 1024
	MBits = 1024 * KBits
	GBits = 1024 * MBits
	TBits = 1024 * GBits

	KBytes = 1024
	MBytes = 1024 * KBytes
	GBytes = 1024 * MBytes
	TBytes = 1024 * GBytes
)

var (
	ErrEmptyArgument = errors.New("empty argument")
)

// 解析2KB，3.13M，4.5G，5.67Tb，字符串到数值
func StringToByte(s string) (uint64, error) {
	if s == "" {
		return 0, ErrEmptyArgument
	}
	i := len(s) - 1
	if s[i] == 'B' || s[i] == 'b' {
		i--
		if i < 0 {
			return 0, fmt.Errorf("invalid argument '%v'", s)
		}
	}
	switch s[i] {
	case 'K', 'k':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * KBytes), nil
	case 'M', 'm':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * MBytes), nil
	case 'G', 'g':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * GBytes), nil
	case 'T', 't':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * TBytes), nil
	default:
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v), nil
	}
}

// 数值1234567，这样的到字符串1.177...MB，保留p位小数
func ByteToString(n uint64, p int) string {
	if n > TBytes {
		return strconv.FormatFloat(float64(n)/TBytes, 'f', p, 64) + "TB"
	}
	if n > GBytes {
		return strconv.FormatFloat(float64(n)/GBytes, 'f', p, 64) + "GB"
	}
	if n > MBytes {
		return strconv.FormatFloat(float64(n)/MBytes, 'f', p, 64) + "MB"
	}
	if n > KBytes {
		return strconv.FormatFloat(float64(n)/KBytes, 'f', p, 64) + "KB"
	}
	return strconv.FormatUint(n, 10) + "B"
}
