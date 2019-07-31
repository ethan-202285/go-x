package weapp_test

import (
	"testing"

	"github.com/goodwong/go-x/wechat/weapp"
)

var (
	dd *weapp.Weapp
)

// 测试之前，请设置Env变量，命令行：
// export CGO_ENABLED=0
// AppID= AppSecret=
// go test ./weapp
// 或
// AppID= AppSecret= go test ./weapp
func init() {
	config := weapp.Config{
		AppID:     "wxc55bbe4cc2995ac5",
		AppSecret: "b493bea4a184cb7a1d783f5c73b499b0",
	}
	dd = weapp.New(&config)
}

func Test_GetAccessToken(t *testing.T) {
	tokenString, err := dd.Client.GetAccessToken()
	if err != nil || len(tokenString) == 0 {
		t.Fatal("GetAccessToken失败：", err)
	}
	t.Log("获取AccessToken成功", tokenString)
}
