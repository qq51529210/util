package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	httpReadSeekerHeaderLen = 1024
)

// HTTPReadSeeker 实现 io.ReadSeeker，
// 服务必须支持 Range 哦
type HTTPReadSeeker struct {
	res *http.Response
	// 资源地址
	url string
	// 请求头的 ETag
	eTag string
	// 当前位置
	offset int64
	// 总字节
	total int64
}

// Close 实现 io.Closer 接口，主要用于关闭 response.body
func (r *HTTPReadSeeker) Close() error {
	if r.res != nil {
		r.res.Body.Close()
		r.res = nil
	}
	return nil
}

// Read 实现 io.Reader 接口
func (r *HTTPReadSeeker) Read(buf []byte) (int, error) {
	n, err := r.res.Body.Read(buf)
	if err != nil {
		return n, err
	}
	r.offset += int64(n)
	//
	return n, nil
}

// Seek 实现 io.Seeker 接口
func (r *HTTPReadSeeker) Seek(offset int64, whence int) (int64, error) {
	// 如果 offset 计算后的距离比 httpReadSeekerHeaderLen 小，
	// 那么没必要再去发起一个请求，直接读就可以了。
	switch whence {
	case io.SeekCurrent:
		// 如果小，读取就行
		if offset < httpReadSeekerHeaderLen {
			_, err := io.CopyN(io.Discard, r.res.Body, offset)
			if err != nil {
				return r.offset, err
			}
			r.offset += offset
			//
			return r.offset, nil
		}
		// 距离很大，重新请求
		offset = r.offset + offset
	case io.SeekEnd:
		// 从开始的偏移
		offset = r.total - offset
		// 走 io.SeekStart 的流程
		fallthrough
	case io.SeekStart:
		// 没变化
		if offset == r.offset {
			return r.offset, nil
		}
		// 在后面
		if offset > r.offset {
			n := offset - r.offset
			// 如果小，读取就行
			if n < httpReadSeekerHeaderLen {
				_, err := io.CopyN(io.Discard, r.res.Body, n)
				if err != nil {
					return r.offset, err
				}
				r.offset += n
				//
				return r.offset, nil
			}
		}
		// 在前面，或者距离很大，重新请求
	default:
		// 错误的
		return r.offset, fmt.Errorf("seek error whence %d", whence)
	}
	//
	err := r.seekHTTP(offset)
	if err != nil {
		return r.offset, err
	}
	return r.offset, nil
}

func (r *HTTPReadSeeker) seekHTTP(offset int64) error {
	req, err := http.NewRequest(http.MethodGet, r.url, nil)
	if err != nil {
		return err
	}
	if r.total < 1 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	} else {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, r.total))
	}
	if r.eTag != "" {
		req.Header.Set("If-Range", r.eTag)
	}
	req.Header.Set("Accept-Ranges", "bytes")
	//
	if r.res != nil {
		r.res.Body.Close()
	}
	r.res, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if r.res.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("error status code %d", r.res.StatusCode)
	}
	r.eTag = r.res.Header.Get("ETag")
	//
	contentRange := r.res.Header.Get("Content-Range")
	if contentRange == "" {
		return errors.New("missing content-range header")
	}
	//
	i := strings.LastIndexByte(contentRange, '/')
	if i < 1 {
		return nil
	}
	n, err := strconv.ParseInt(contentRange[i+1:], 10, 64)
	if err != nil {
		return fmt.Errorf("error content range %s", contentRange)
	}
	r.offset = offset
	r.total = n
	//
	return nil
}

// NewHTTPReadSeeker 返回新的 HTTPReadSeeker ，开始会请求一次
func NewHTTPReadSeeker(url string) (*HTTPReadSeeker, error) {
	r := new(HTTPReadSeeker)
	r.url = url
	err := r.seekHTTP(0)
	if err != nil {
		return nil, err
	}
	return r, nil
}
