package dingtalk_test

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/goodwong/go-x/dingtalk"
)

var (
	dd *dingtalk.Dingtalk
)

// 测试之前，请设置Env变量，命令行：
// export CGO_ENABLED=0
// export CorpID= AgentID= AppKey= AppSecret=
// go test ./dingtalk
// 或
// CorpID= AgentID= AppKey= AppSecret= go test ./dingtalk
func init() {
	agentID, _ := strconv.ParseUint(os.Getenv("AgentID"), 10, 64)
	config := dingtalk.Config{
		CorpID:    os.Getenv("CorpID"),
		AgentID:   agentID,
		AppKey:    os.Getenv("AppKey"),
		AppSecret: os.Getenv("AppSecret"),
	}
	if config.AppKey == "" {
		log.Fatal("\033[31m读取配置失败。请按以下步骤进行测试：\033[0m\n" +
			"\t\033[33m1. 请在钉钉应用管理后台设置IP白名单.\n" +
			"\t2. 设置测试参数并运行：\n" +
			"\t\033[7mCorpID= AgentID= AppKey= AppSecret=\033[0m go test ./dingtalk",
		)
	}
	dd = dingtalk.New(&config)
}

func Test_GetAccessToken(t *testing.T) {
	tokenString, err := dd.Client.GetAccessToken()
	if err != nil || len(tokenString) == 0 {
		t.Fatal("GetAccessToken失败：", err)
	}
	t.Log("获取AccessToken成功", tokenString)
}

func Test_SendWorkText(t *testing.T) {
	taskID, err := dd.SendText(map[string]interface{}{"userid_list": "manager7140"}, "一个小小的开始，不要骄傲！")
	if err != nil {
		t.Fatal("dd.SendText()失败：", err)
	}
	t.Log("dd.SendText()成功", taskID)
}

func Test_LoginByCode(t *testing.T) {
	info, err := dd.UserInfoByCode("8234612394192634916234")
	if !strings.Contains(err.Error(), "不存在的临时授权码") {
		t.Fatal("dd.LoginByCode()失败：", err)
	}
	t.Log("dd.UserInfoByCode()测试通过", info)
}
