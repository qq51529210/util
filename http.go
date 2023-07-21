package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

// HTTPError 表示 JSON 错误
type HTTPError struct {
	Phrase string
	Detail string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("phrase %s, detail %s", e.Phrase, e.Detail)
}

// HTTPStatusError 表示状态错误
type HTTPStatusError int

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("status code %d", e)
}

// HTTP 封装 http 操作
// method 方法
// url 请求地址
// query 请求参数
// reqBody 用于读取发送 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
// timeout 超时
func HTTP[reqData, resData any](method, url string, query url.Values, reqBody *reqData, resBody *resData, onResponse func(res *http.Response) error, timeout time.Duration) error {
	// 请求
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if onResponse != nil {
		err = onResponse(res)
		if err != nil {
			return err
		}
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return nil
}

// HTTPTo 封装 http 操作
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// resBody 写入响应的 body 数据
// statusCode 用于判断状态码
// timeout 超时
func HTTPTo[reqData any](method, url string, query url.Values, reqBody *reqData, resBody io.Writer, onResponse func(res *http.Response) error, timeout time.Duration) error {
	// 请求
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()u
	// 状态码
	if onResponse != nil {
		err = onResponse(res)
		if err != nil {
			return err
		}
	}
	// 解析
	if resBody != nil {
		_, err = io.Copy(resBody, res.Body)
	}
	return err
}

// HTTPFrom 封装 http 操作
// method 方法
// url 请求地址
// query 请求参数
// reqBody 用于读取发送 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
// timeout 超时
func HTTPFrom[resData any](method, url string, query url.Values, reqBody io.Reader, resBody *resData, onResponse func(res *http.Response) error, timeout time.Duration) error {
	// 请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if onResponse != nil {
		err = onResponse(res)
		if err != nil {
			return err
		}
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return err
}

// HTTPWithContext 封装 http 操作
// ctx 超时上下文
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
func HTTPWithContext[reqData, resData any](ctx context.Context, method, url string, query url.Values, reqBody *reqData, resBody *resData, onResponse func(res *http.Response) error) error {
	// body
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	// 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// query
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if onResponse != nil {
		err = onResponse(res)
		if err != nil {
			return err
		}
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return nil
}

// HTTPToWithContext 封装 http 操作
// ctx 超时上下文
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// resBody 写入响应的 body 数据
// statusCode 用于判断状态码
func HTTPToWithContext[reqData any](ctx context.Context, method, url string, query url.Values, reqBody *reqData, resBody io.Writer, onResponse func(res *http.Response) error) error {
	// body
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	// 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// query
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if onResponse != nil {
		err = onResponse(res)
		if err != nil {
			return err
		}
	}
	// 解析
	if resBody != nil {
		_, err = io.Copy(resBody, res.Body)
	}
	return err
}

// HTTPFromWithContext 封装 http 操作
// ctx 超时上下文
// method 方法
// url 请求地址
// query 请求参数
// reqBody 用于读取发送 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
func HTTPFromWithContext[resData any](ctx context.Context, method, url string, query url.Values, reqBody io.Reader, resBody *resData, onResponse func(res *http.Response) error) error {
	// 请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}
	// query
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if onResponse != nil {
		err = onResponse(res)
		if err != nil {
			return err
		}
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return err
}

var (
	// HTTPQueryTag 是 HTTPQuery 解析 tag 的名称
	HTTPQueryTag = "query"
)

// HTTPQuery 将结构体 v 格式化到 url.Values
// 只扫描一层，并略过空值
func HTTPQuery(v any, q url.Values) url.Values {
	if q == nil {
		q = make(url.Values)
	}
	rv := reflect.ValueOf(v)
	vk := rv.Kind()
	if vk == reflect.Pointer {
		rv = rv.Elem()
		vk = rv.Kind()
	}
	if vk != reflect.Struct {
		panic("v must be struct or struct ptr")
	}
	return httpQuery(rv, q)
}

func httpQuery(v reflect.Value, q url.Values) url.Values {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构，只一层
		if fvk == reflect.Struct {
			continue
		}
		if fvk == reflect.String {
			// 空字符串
			if fv.IsZero() {
				continue
			}
		}
		ft := vt.Field(i)
		tn := ft.Tag.Get(HTTPQueryTag)
		if tn == "" {
			continue
		}
		q.Set(tn, fmt.Sprintf("%v", fv.Interface()))
	}
	return q
}
