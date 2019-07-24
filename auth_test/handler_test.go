package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goodwong/go-x/auth"
)

var (
	testHandlerTokens *auth.TokenResponse
)

func TestLoginMissParams(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost/api/login", nil)
	w := httptest.NewRecorder()
	auths.Handler.HandleLogin(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	if !strings.Contains(content, "invalid provider") {
		t.Fatalf("非期望的响应：%s", content)
	}
}

func TestLoginByPassword(t *testing.T) {
	// POST参数
	buffer := bytes.NewBufferString(`{
		"username": "testpassword",
		"password": "testpassWord123,"
	}`)

	// 模拟请求
	req := httptest.NewRequest("POST", "http://localhost/api/login?provider=password&remember=1&device=gotest", buffer)
	w := httptest.NewRecorder()
	auths.Handler.HandleLogin(w, req)

	// 测试结果
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("非200响应：%d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	t.Log(content)
	if !strings.Contains(content, `"token_expires"`) {
		t.Fatalf("非期望的响应：%s", content)
	}

	// 解析结果，后面要用
	tokens := &auth.TokenResponse{}
	err := json.Unmarshal(body, tokens)
	if err != nil {
		t.Fatal(err)
	}
	testHandlerTokens = tokens
}

func TestRenew(t *testing.T) {
	// POST参数
	buffer := bytes.NewBufferString(fmt.Sprintf(`{"refresh_token": "%s"}`, *testHandlerTokens.RefreshToken))

	// 模拟请求
	req := httptest.NewRequest("POST", "http://localhost/api/login?provider=password", buffer)
	w := httptest.NewRecorder()
	auths.Handler.HandleRenew(w, req)

	// 测试结果
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("非200响应：%d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	if !strings.Contains(content, `"token_expires"`) {
		t.Fatalf("非期望的响应：%s", content)
	}
}

// 测试非登录态时候的登出行为
func TestLogout(t *testing.T) {
	// 模拟请求
	req := httptest.NewRequest("DELETE", "http://localhost/api/login", nil)
	w := httptest.NewRecorder()
	auths.Handler.HandleLogout(w, req)

	// 测试结果
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	if !strings.Contains(content, `无需登出`) {
		t.Fatalf("非期望的响应：%s", content)
	}
}

// 测试登录态时候的登出行为
func TestLogout2(t *testing.T) {
	// 模拟请求
	// 加上jwt令牌
	url := fmt.Sprintf("http://localhost/api/login?jwt=%s", testHandlerTokens.Token)
	req := httptest.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	// 加上middleware
	handlerFunc := http.HandlerFunc(auths.Handler.HandleLogout)
	handler := auths.Middleware.AuthenticatedWithUser(handlerFunc)
	handler.ServeHTTP(w, req)

	// 测试结果
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("非200响应：%d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	if !strings.Contains(content, `登出成功`) {
		t.Fatalf("非期望的响应：%s", content)
	}
}
