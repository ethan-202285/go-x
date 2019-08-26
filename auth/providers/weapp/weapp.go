package weapp

import (
	"encoding/json"

	"github.com/goodwong/go-x/auth"
	"github.com/goodwong/go-x/wechat/weapp"
	"github.com/jinzhu/gorm"
)

// NewProvider 创建实例
func NewProvider(config *Config) *Provider {
	return &Provider{
		auth:  config.Auth,
		weapp: config.Weapp,
	}
}

// Config 配置
type Config struct {
	Auth  *auth.Auth
	Weapp *weapp.Weapp
}

// Provider 通过dingtalk 登陆
type Provider struct {
	auth  *auth.Auth
	weapp *weapp.Weapp
}

func (p *Provider) repository() *auth.Repository {
	return p.auth.Repository
}

// Name 获取provider名字; implemented Name with LoginProvider interface
func (p *Provider) Name() string {
	return "wechat_weapp"
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
	// 接口获取数据
	session, err := p.weapp.Code2session(credentials.Code)
	if err != nil {
		return nil, err
	}

	openID := session.OpenID + "@" + p.weapp.AppID
	unionID := ""
	if session.UnionID != "" {
		unionID = session.UnionID
	}

	// 如果用户存在，直接返回
	// SELECT * FROM "user_identities" WHERE (provider = 'wechat_weapp' and open_id = 'oPy_U5****Q6y7so@wx70****f05') LIMIT 1
	user, err = p.repository().FindByOpenID(p.Name(), openID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// 创建用户
	// 优先以unionID 创建用户
	var username string
	if unionID != "" {
		username = unionID + "@" + "wechat_unionID"
	} else {
		username = openID + "@" + p.Name()
	}
	// 如果用户已存在，直接关联该用户
	user, _ = p.repository().FindByUsername(username)
	if user == nil {
		// 否则就要创建用户
		user, err = p.repository().Create(username, "")
		if err != nil {
			return nil, err
		}
	}
	// 然后创建登陆凭证
	_, err = p.repository().CreateIdentity(user.ID, p.Name(), openID)
	if err != nil {
		return nil, err
	}
	return user, err
}
