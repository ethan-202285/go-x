package auth

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// NewPasswordProvider 创建实例
// Usage:
// ```go
//   passwordProvider := NewPasswordProvider(auth)
//   auth.Service.RegisterProvider(passwordProvider)
// ```
func NewPasswordProvider(auth *Auth) *PasswordProvider {
	return &PasswordProvider{auth: auth}
}

// PasswordProvider 通过password 登陆
type PasswordProvider struct {
	auth *Auth
}

func (p *PasswordProvider) repository() *Repository {
	return p.auth.Repository
}

// Name 获取provider名字; implemented Name with LoginProvider interface
func (p *PasswordProvider) Name() string {
	return "password"
}

func (p *PasswordProvider) passwordMatched(identity *UserIdentity, password string) bool {
	data := struct {
		PasswordHash string `json:"password_hash"`
	}{}
	if err := json.Unmarshal([]byte(*identity.Data), &data); err != nil {
		return false
	}
	// 打包
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []byte(password))
	binary.Write(buf, binary.BigEndian, p.auth.secretKey)

	err := bcrypt.CompareHashAndPassword([]byte(data.PasswordHash), buf.Bytes())
	return err == nil
}

// Login 登陆; implemented Login with LoginProvider interface
func (p *PasswordProvider) Login(payload []byte) (user *User, err error) {
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
func (p *PasswordProvider) passwordHash(password string) string {
	// 打包
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []byte(password))
	binary.Write(buf, binary.BigEndian, p.auth.secretKey)
	hash, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost) // 50ms
	if err != nil {
		panic(err)
	}

	// base64
	return string(hash) // bcrypt 本身有base64编码
}

// Register 注册用户
func (p *PasswordProvider) Register(username, password string, userID ...uint64) (user *User, err error) {
	// 如果指定用户
	if len(userID) == 1 {
		user, err = p.repository().Find(userID[0])
	}
	if err != nil {
		return nil, err
	}

	// 创建用户
	user, err = p.repository().Create(username, username)
	if err != nil {
		return nil, err
	}

	// 创建密码
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
