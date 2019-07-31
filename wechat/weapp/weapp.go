package weapp

import (
	"strings"

	"github.com/goodwong/go-x/dingtalk/client"
)

// New 创建
func New(cfg *Config) *Weapp {
	r := strings.NewReplacer("APPID", cfg.AppID, "APPSECRET", cfg.AppSecret)
	tokenAPI := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET"
	tokenAPI = r.Replace(tokenAPI)

	client := client.New(&client.Config{
		TokenAPI: tokenAPI,
	})
	return &Weapp{
		config: cfg,
		Client: client,
	}
}

// Config 配置类
type Config struct {
	AppID     string
	AppSecret string
}

// Weapp 功能类
type Weapp struct {
	config *Config
	Client *client.APIClient
}
