package dingtalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	// AccessTokenFunc 刷新access_token的方法，
	// 默认直接从钉钉接口获取；
	// 也可以从Redis获取（中控服务器会统一管理access_token）
	AccessTokenFunc = defaultAccessTokenFunc
)

type accessTokenStruct struct {
	Token     string
	ExpiredAt time.Time
}

// defaultAccessTokenFunc 默认的获取access_token的具体方法
func defaultAccessTokenFunc(cfg *Config) (accessToken *accessTokenStruct, err error) {
	r := strings.NewReplacer("KEY", cfg.AppKey, "SECRET", cfg.AppSecret)
	api := r.Replace("https://oapi.dingtalk.com/gettoken?appkey=KEY&appsecret=SECRET")

	retryFn := func() (result interface{}, retry bool, err error) {
		// 网络原因，重试！
		resp, err := newHTTPClient().Get(api)
		if err != nil {
			return nil, true, fmt.Errorf("请求失败: %s", err)
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, true, fmt.Errorf("读取失败: %s", err)
		}

		// 返回数据：
		contentType := resp.Header.Get("Content-type")
		status := resp.StatusCode
		errInfo, err := checkResult(status, contentType, respBytes)
		// 服务器忙碌，重试！
		if errInfo.Errcode == -1 {
			return nil, true, errors.New(errInfo.Errmsg)
		}

		// 其它错误，返回
		if err != nil {
			return nil, false, fmt.Errorf("钉钉返回错误：%s", err)
		}

		// 解析
		jsonResult := &struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
		}{}
		err = json.Unmarshal(respBytes, jsonResult)
		if err != nil {
			return nil, false, fmt.Errorf("解析返回的JSON错误：%s", err)
		}
		// 一切正常：返回
		return &accessTokenStruct{
			Token:     jsonResult.AccessToken,
			ExpiredAt: time.Now().Add(time.Duration(7200-120) * time.Second),
		}, false, nil
	}
	result, err := retry(retryFn, 3) // 重试3次
	if err != nil {
		return nil, err
	}
	accessToken, _ = result.(*accessTokenStruct)
	return accessToken, nil
}

// GetAccessToken 获取access_token（可能是缓存）
func (dd *Dingtalk) GetAccessToken() (tokenString string, err error) {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	accessToken := dd.accessToken
	if accessToken != nil && accessToken.ExpiredAt.After(time.Now()) {
		return accessToken.Token, nil
	}
	return dd.fetchToken()
}

// RefreshAccessToken 强制从服务器刷新access_token
func (dd *Dingtalk) RefreshAccessToken() (tokenString string, err error) {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	return dd.fetchToken()
}

func (dd *Dingtalk) fetchToken() (tokenString string, err error) {
	accessToken, e := AccessTokenFunc(dd.config)
	// 失败
	if e != nil {
		return "", e
	}
	// 成功
	dd.accessToken = accessToken // 如果出错，返回的是nil
	log.Println("获取AccessToken成功")
	return accessToken.Token, nil
}
