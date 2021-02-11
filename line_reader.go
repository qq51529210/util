package common

import (
	"bytes"
	"io"
)

// 读取一行数据
type LineReader struct {
	buff   []byte    // 读取数据的缓存
	data   []byte    // 数据
	pIdx   int       // data解析的索引
	dIdx   int       // data有效数据的索引
	dLen   int       // data的有效数据大小
	reader io.Reader // 数据源
}

// 创建一个LineReader，数据源是r，r.Read()使用的数据缓存大小是n。
func NewLineReader(r io.Reader, n int) *LineReader {
	p := new(LineReader)
	if n <= 0 {
		n = 128
	}
	p.buff = make([]byte, n)
	p.reader = r
	return p
}

// 重新设置数据源r
func (r *LineReader) SetReader(reader io.Reader) {
	r.reader = reader
}

// 读取一行数据
func (r *LineReader) ReadLine() ([]byte, error) {
	// 从缓存中读取数据
	p := r.readData()
	if p != nil {
		return p, nil
	}
	// 从reader读取数据
	var n int
	var err error
	for {
		// 从reader读取数据
		n, err = r.reader.Read(r.buff)
		if err != nil {
			// 没有数据了
			if err == io.EOF {
				// 返回data
				p = r.data[r.dIdx:r.dLen]
				r.dIdx = 0
				r.pIdx = 0
				r.dLen = 0
			}
			return nil, err
		}
		// data没有数据，先解析buff，减小拷贝
		if r.dLen == 0 {
			// buff中是否有完整的一行
			i := bytes.IndexByte(r.buff[:n], '\n')
			if i >= 0 {
				p = r.checkEnter(r.buff[:i])
				if p != nil {
					// 添加buff[m:n]到data
					r.appendData(r.buff[i+1 : n])
					// 返回
					return p, nil
				}
			}
			// 添加buff[:n]到data
			r.appendData(r.buff[:n])
			// 继续读取
			continue
		}
		// data中有数据，添加buff[:n]到data
		r.appendData(r.buff[:n])
		// 从data中读取数据
		p := r.readData()
		if p != nil {
			return p, nil
		}
	}
}

// 从缓存中读取数据
func (r *LineReader) readData() []byte {
	// 有数据
	if r.dLen > r.pIdx {
		i := bytes.IndexByte(r.data[r.pIdx:r.dLen], '\n')
		if i >= 0 {
			r.pIdx = r.pIdx + i
			p := r.checkEnter(r.data[r.dIdx:r.pIdx])
			r.pIdx++
			if r.pIdx >= r.dLen {
				r.pIdx = 0
				r.dLen = 0
			}
			return p
		}
		r.pIdx = r.dLen
	}
	return nil
}

// 添加数据到data缓存
func (r *LineReader) appendData(b []byte) {
	i := copy(r.data[r.dLen:], b)
	if i < len(b) {
		r.data = append(r.data, b[i:]...)
	}
	r.dLen += len(b)
}

// 检查最后一个字符是否'\r'
func (r *LineReader) checkEnter(p []byte) []byte {
	i := len(p) - 1
	if i >= 0 && p[i] == '\r' {
		return p[:i]
	}
	return p
}
