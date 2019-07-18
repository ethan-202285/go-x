package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// DefaultTokenLife 默认JWT token有效时长
var DefaultTokenLife = 1 * time.Hour // 1 Hour

func newService(auth *Auth) *Service {
	return &Service{auth: auth, providers: map[string]LoginProvider{}}
}

// LoginProvider 登陆方法
type LoginProvider interface {
	Name() string
	Login(credentials []byte) (user *User, err error)
}

// Service Auth 服务
type Service struct {
	auth      *Auth
	providers map[string]LoginProvider
}

func (s *Service) repository() *Repository {
	return s.auth.Repository
}

// issueRefreshToken 获取用于refresh的token
func (s *Service) issueRefreshToken(
	user *User, device, remark string,
) (tokenString string, expires time.Time, err error) {
	// 同一个设备只能有一个有效token，创建新token，该device的原有token失效
	if err = s.repository().DeleteToken(user.ID, device); err != nil {
		return
	}
	// 颁发新token，返回的tokenString返回到前端保存
	token, tokenString, err := s.repository().CreateToken(user.ID, device, remark)
	if token != nil {
		expires = token.ExpiredAt
	}
	return
}

// issueJWTToken 获取jwt的token
func (s *Service) issueJWTToken(
	user *User,
) (tokenString string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(DefaultTokenLife)
	claims := jwt.MapClaims{
		"iat": now.UTC().Unix(),
		"sub": user.ID,
		"exp": expires.UTC().Unix(),
	}
	_, tokenString, err = s.auth.jwtauth.Encode(claims)
	return
}

// TokenResponse 返回token结构
type TokenResponse struct {
	Token               string     `json:"token"`
	TokenExpires        time.Time  `json:"token_expires"`
	RefreshToken        *string    `json:"refresh_token,omitempty"`
	RefreshTokenExpires *time.Time `json:"refresh_token_expires,omitempty"`
}

// Login 登陆
func (s *Service) Login(
	providerName string, credentials []byte, remember bool, device string, deviceName ...string,
) (tokens *TokenResponse, err error) {
	// login
	provider, ok := s.providers[providerName]
	if !ok {
		return nil, errors.New("invalid provider")
	}
	user, err := provider.Login(credentials)
	if err != nil {
		return nil, err
	}

	// issue tokens
	tokens = &TokenResponse{}
	if remember {
		remark := ""
		if len(deviceName) == 1 {
			remark = deviceName[0]
		}
		// 如果出错也没关系，重点是jwt要成功
		refreshToken, refreshTokenExpires, err := s.issueRefreshToken(user, device, remark)
		if err == nil {
			tokens.RefreshToken = &refreshToken
			tokens.RefreshTokenExpires = &refreshTokenExpires
		}
	}
	tokens.Token, tokens.TokenExpires, err = s.issueJWTToken(user)
	if err != nil {
		return nil, err
	}
	return
}

// Renew 通过RefreshToken续约
func (s *Service) Renew(tokenString string) (tokens *TokenResponse, err error) {
	// 验证
	user, err := s.repository().FindByToken(tokenString)
	if err != nil {
		return nil, err
	}
	// 发放jwttoken
	tokens = &TokenResponse{}
	tokens.Token, tokens.TokenExpires, err = s.issueJWTToken(user)
	if err != nil {
		return nil, err
	}
	return
}

// Logout 登出
// todo 需要加入 Logouts 名单
// todo middleware 检查这个Logouts名单
func (s *Service) Logout(user *User, device string) (err error) {
	return s.repository().DeleteToken(user.ID, device)
}

// RegisterProvider 注册登陆方式
func (s *Service) RegisterProvider(provider LoginProvider) {
	name := provider.Name()
	s.providers[name] = provider
	return
}
