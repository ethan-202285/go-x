package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const defaultHTTPClientTimeout = time.Duration(2 * time.Second)

// 封装的request
// 做两件事：
// 1. 自动替换ACCESS_TOKEN
// 2. ACCESS_TOKEN过期自动重试
func (dd *Dingtalk) request(method, url, contentType string, body io.Reader) (respBytes []byte, err error) {
	retryFn := func() (result interface{}, retry bool, err error) {
		// 替换access_token
		token, err := dd.GetAccessToken()
		if err != nil {
			// GetAccessToken 本身会遇忙重试
			return nil, false, fmt.Errorf("获取AccessToken失败：%s", err)
		}
		r := strings.NewReplacer("ACCESS_TOKEN", token)

		// request
		client := newHTTPClient()
		req, err := http.NewRequest(method, r.Replace(url), body)
		if err != nil {
			return nil, true, err
		}
		req.Header.Set("Content-Type", contentType)
		var resp *http.Response
		resp, err = client.Do(req)
		if err != nil {
			return nil, true, fmt.Errorf("请求失败：%s", err)
		}

		// response
		respBytes, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, true, fmt.Errorf("读取失败：%s", err)
		}

		// 验证
		contentType := resp.Header.Get("Content-type")
		status := resp.StatusCode
		var errInfo responseErrorInfo
		errInfo, err = checkResult(status, contentType, respBytes)
		// 需要重试的情况
		switch errInfo.Errcode {
		case -1:
			return nil, true, errors.New(errInfo.Errmsg) // 服务器忙
		case 40001, 40014, 41001, 42001:
			log.Println("access_token失效！", err)
			if _, err = dd.RefreshAccessToken(); err != nil { // accessToken 失效，重新获取，并重试
				return nil, false, fmt.Errorf("刷新access_token失败：%s", err)
			}
			return nil, true, fmt.Errorf("access_token失效：%s", err)
		}
		// 其它错误，返回
		if err != nil {
			return nil, false, fmt.Errorf("钉钉返回错误：%s", err)
		}

		// 正常：返回
		return respBytes, false, nil
	}
	result, err := retry(retryFn, 3) // 重试3次
	if err != nil {
		return nil, err
	}
	respBytes, _ = result.([]byte)
	return respBytes, nil
}

// Get 请求
func (dd *Dingtalk) Get(url string) (respBytes []byte, err error) {
	return dd.request("GET", url, "application/json", nil)
}

// PostJSON 请求
func (dd *Dingtalk) PostJSON(url string, data interface{}) (respBytes []byte, err error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false) // 禁用转码
	_ = encoder.Encode(data)

	return dd.request("POST", url, "application/json", buffer)
}

// PostBlob 请求
func (dd *Dingtalk) PostBlob(url, field string, blob []byte) (respBytes []byte, err error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile(field, "media.png")
	if err != nil {
		return
	}
	part.Write(blob)
	err = writer.Close()
	if err != nil {
		return
	}

	return dd.request("POST", url, writer.FormDataContentType(), buffer)
}

type retryFn func() (result interface{}, retry bool, err error)

func retry(fn retryFn, retryTimes int) (result interface{}, err error) {
	retry := 0
	for {
		result, tryAgain, err := fn()

		// 成功
		if err == nil {
			return result, nil
		}

		// 失败，重试
		if retry < retryTimes && tryAgain {
			retry++
			continue
		}

		// 重试机会用尽，或无需重试
		return nil, err
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultHTTPClientTimeout,
	}
}

type responseErrorInfo struct {
	Errcode int
	Errmsg  string
}

func checkResult(status int, contentType string, readBytes []byte) (responseErrorInfo, error) {
	var err error
	var result responseErrorInfo
	if strings.Contains(contentType, "application/json") {
		// json解码
		err = json.Unmarshal(readBytes, &result)
		if err != nil {
			return result, err
		}
		// 业务检查
		if result.Errcode != 0 {
			return result, errors.New(result.Errmsg)
		}
	}
	// 检查status code
	if (status < 200 || status >= 300) && err == nil {
		err = fmt.Errorf("http状态码错误：%d", status)
	}
	return result, err
}
