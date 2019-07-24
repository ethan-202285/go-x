package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// DefaultTokenLife 默认JWT token有效时长
var DefaultTokenLife = 1 * time.Hour // 1 Hour
var cleanupInterval = 10 * time.Second

func newService(auth *Auth) *Service {
	service := &Service{auth: auth, providers: map[string]LoginProvider{}}
	service.cleanupLogoutsLoop()
	return service
}

// Service Auth 服务
type Service struct {
	auth      *Auth
	providers map[string]LoginProvider
	logouts   sync.Map
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
func (s *Service) Renew(tokenString string) (user *User, tokens *TokenResponse, err error) {
	// 验证
	user, err = s.repository().FindByToken(tokenString)
	if err != nil {
		return nil, nil, err
	}
	// 发放jwttoken
	tokens = &TokenResponse{}
	tokens.Token, tokens.TokenExpires, err = s.issueJWTToken(user)
	if err != nil {
		return nil, nil, err
	}
	return
}

// Logout 登出
// todo middleware 检查这个Logouts名单
func (s *Service) Logout(user *User, device string) (err error) {
	logout := logoutStruct{UserID: user.ID, LogoutAt: time.Now()}
	s.logouts.Store(user.ID, &logout)
	return s.repository().DeleteToken(user.ID, device)
}

// JwtInvalid 检查是否jwt是否提前失效（指用户主动登出）
func (s *Service) JwtInvalid(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)
	userID := uint64(claims["sub"].(float64))
	issuedAt := int64(claims["iat"].(float64))

	// 查询
	// 如果没在logouts里面找到，
	// 说明用户没有主动注销行为，jwt可以继续使用
	v, ok := s.logouts.Load(userID)
	if !ok {
		return false
	}

	// 比较
	// 如果是注销前颁发的jwt，则失效，需要重新登录
	logout := v.(*logoutStruct)
	if issuedAt <= logout.LogoutAt.UTC().Unix() {
		return true
	}

	// 注销后颁发的jwt，可以继续使用那个
	return false
}

// logoutStruct 主动注销行为（数据量较小，可内存保存）
//（只要是userID在这里的，并且user.iat < LogoutAt的，都要重新登录）
type logoutStruct struct {
	UserID   uint64
	LogoutAt time.Time
}

// 自动清理logouts条目
// todo 集成CanceledContext，感知cancel时间
func (s *Service) cleanupLogoutsLoop() {
	go func() {
		for {
			now := time.Now()
			s.logouts.Range(func(k, v interface{}) bool {
				logout := v.(*logoutStruct)
				if logout.LogoutAt.Add(DefaultTokenLife).Before(now) {
					s.logouts.Delete(k)
				}

				return true
			})

			// 间隔
			time.Sleep(cleanupInterval)
		}
	}()
}
