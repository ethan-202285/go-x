package dingtalk

import (
	"encoding/json"

	"github.com/goodwong/go-x/auth"
	"github.com/goodwong/go-x/dingtalk"
	"github.com/jinzhu/gorm"
)

// NewProvider 创建实例
func NewProvider(config *Config) *Provider {
	return &Provider{
		auth:     config.Auth,
		dingtalk: config.Dingtalk,
	}
}

// Config 配置
type Config struct {
	Auth     *auth.Auth
	Dingtalk *dingtalk.Dingtalk
}

// Provider 通过dingtalk 登陆
type Provider struct {
	auth     *auth.Auth
	dingtalk *dingtalk.Dingtalk
}

func (p *Provider) repository() *auth.Repository {
	return p.auth.Repository
}

// Name 获取provider名字; implemented Name with LoginProvider interface
func (p *Provider) Name() string {
	return "dingtalk"
}

// Login 登陆; implemented Login with LoginProvider interface
func (p *Provider) Login(payload []byte) (user *auth.User, err error) {
	// params
	credentials := struct {
		Code string
	}{}
	if err := json.Unmarshal(payload, &credentials); err != nil {
		return nil, err
	}

	// 钉钉接口获取数据
	info, err := p.dingtalk.UserInfoByCode(credentials.Code)
	if err != nil {
		return nil, err
	}

	// 如果用户存在，直接返回
	providerName := p.Name()
	user, err = p.repository().FindByOpenID(providerName, info.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// 用户不存在，说明这个ding user没有登陆过
	// 此时需要再次检查这个手机号码的用户是否已经存在了
	// 如果不存在，创建用户
	username := info.Mobile + "@telephone"
	user, err = p.repository().FindByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		user, err = p.repository().Create(username, info.Name, info.Avatar)
		if err != nil {
			return nil, err
		}
	}
	// 然后创建登陆凭证，并关联至这个手机号码用户
	_, err = p.repository().CreateIdentity(user.ID, p.Name(), info.UserID, info)
	if err != nil {
		return nil, err
	}
	return user, err
}
