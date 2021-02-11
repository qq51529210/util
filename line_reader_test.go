package common

import (
	"bytes"
	"testing"
)

func TestLienReader(t *testing.T) {
	var data bytes.Buffer
	reader := NewLineReader(&data, 32)
	var line []byte
	// 读取空行
	data.WriteByte('\n')
	line, _ = reader.ReadLine()
	if len(line) != 0 {
		t.FailNow()
	}
	data.Reset()
	data.WriteByte('\r')
	data.WriteByte('\n')
	line, _ = reader.ReadLine()
	if len(line) != 0 {
		t.FailNow()
	}
	// 读取一行
	var lines []string
	lines = append(lines, "asdfsdf")
	lines = append(lines, "wweqgfdfbv")
	lines = append(lines, "xcvxcv")
	lines = append(lines, "xcgkfjghjvxcv")
	data.Reset()
	for i, line := range lines {
		data.WriteString(line)
		if i%2 == 0 {
			data.WriteByte('\n')
		} else {
			data.WriteString("\r\n")
		}
	}
	// 不写入换行
	lines = append(lines, "234324123")
	data.WriteString(lines[len(lines)-1])
	// 检查
	i := 0
	for {
		line, _ = reader.ReadLine()
		if line == nil {
			break
		}
		s := string(line)
		if s != lines[i] {
			t.FailNow()
		}
		i++
	}
}
