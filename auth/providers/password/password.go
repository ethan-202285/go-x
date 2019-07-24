package password

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"unicode"

	"github.com/goodwong/go-x/auth"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// NewProvider 创建实例
// Usage:
// ```go
//   provider := NewProvider(auth)
//   auth.Service.RegisterProvider(provider)
// ```
func NewProvider(config *Config) *Provider {
	return &Provider{
		auth:      config.Auth,
		secretKey: config.SecretKey,
	}
}

// Config 配置
type Config struct {
	Auth      *auth.Auth
	SecretKey []byte
}

// Provider 通过password 登陆
type Provider struct {
	auth      *auth.Auth
	secretKey []byte
}

func (p *Provider) repository() *auth.Repository {
	return p.auth.Repository
}

// Name 获取provider名字; implemented Name with LoginProvider interface
func (p *Provider) Name() string {
	return "password"
}

func (p *Provider) passwordMatched(identity *auth.UserIdentity, password string) bool {
	data := struct {
		PasswordHash string `json:"password_hash"`
	}{}
	if err := json.Unmarshal([]byte(*identity.Data), &data); err != nil {
		return false
	}
	// 打包
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []byte(password))
	binary.Write(buf, binary.BigEndian, p.secretKey)

	err := bcrypt.CompareHashAndPassword([]byte(data.PasswordHash), buf.Bytes())
	return err == nil
}

// Login 登陆; implemented Login with LoginProvider interface
func (p *Provider) Login(payload []byte) (user *auth.User, err error) {
	// params
	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := json.Unmarshal(payload, &credentials); err != nil {
		return nil, err
	}

	// validate password
	providerName := p.Name()
	identity, err := p.repository().FindIdentity(providerName, credentials.Username)
	if err != nil {
		return nil, errors.New("无效的用户名或密码")
	}
	if !p.passwordMatched(identity, credentials.Password) {
		return nil, errors.New("无效的用户名或密码")
	}

	// find user
	user, err = p.repository().Find(identity.UserID)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("无效的用户名或密码")
	}
	return
}

// PasswordHash 获取hash的密码
func (p *Provider) passwordHash(password string) string {
	// 打包
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []byte(password))
	binary.Write(buf, binary.BigEndian, p.secretKey)
	hash, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost) // 50ms
	if err != nil {
		panic(err)
	}

	// base64
	return string(hash) // bcrypt 本身有base64编码
}

// Register 注册用户
func (p *Provider) Register(username, password string, bindUserID ...uint64) (user *auth.User, err error) {
	if len(bindUserID) == 1 {
		// 如果指定用户
		user, err = p.repository().Find(bindUserID[0])
	} else {
		// 创建用户
		user, err = p.repository().Create(username, username)
	}
	if err != nil {
		return nil, err
	}

	// 创建密码
	if !isValid(password) {
		return nil, errors.New("密码长度须大于8位，包含大小写，特殊字符、数字")
	}
	passwordHash := p.passwordHash(password)
	data := struct {
		PasswordHash string `json:"password_hash"`
	}{
		PasswordHash: passwordHash,
	}
	_, err = p.repository().CreateIdentity(user.ID, p.Name(), username, data)
	if err != nil {
		return nil, err
	}
	return
}

// SetPassword 重设密码
func (p *Provider) SetPassword(username, password string) (err error) {
	// 找到用户
	identity, err := p.repository().FindIdentity(p.Name(), username)
	if err != nil {
		return err
	}

	// 创建密码
	if !isValid(password) {
		return errors.New("密码长度须大于8位，包含大小写，特殊字符、数字")
	}
	passwordHash := p.passwordHash(password)
	data := struct {
		PasswordHash string `json:"password_hash"`
	}{
		PasswordHash: passwordHash,
	}

	// 设置密码
	p.repository().UpdateIdentityData(identity, data)
	return nil
}

// https://stackoverflow.com/a/56139457
func isValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) > 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
