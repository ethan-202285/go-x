package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
)

// go get -u github.com/kevinburke/go-bindata/...
//go:generate $GOPATH/bin/go-bindata -pkg=auth template/...

// New 返回Auth类
func New(config Config) *Auth {
	auth := &Auth{
		secretKey: config.SecretKey,
		gormDB:    config.DB,
		jwtauth:   jwtauth.New("HS256", config.SecretKey, nil),
	}
	auth.Repository = newRepository(auth)
	auth.Service = newService(auth)
	auth.Handler = newHandler(auth)
	auth.Middleware = newMiddleware(auth)
	return auth
}

// Config 配置
type Config struct {
	SecretKey []byte
	DB        *gorm.DB
}

// Auth 认证类
type Auth struct {
	Repository *Repository
	Service    *Service
	Handler    *Handler
	Middleware *Middleware

	secretKey []byte
	gormDB    *gorm.DB
	jwtauth   *jwtauth.JWTAuth
}

func (auth *Auth) db() *gorm.DB {
	if auth.gormDB == nil {
		panic("Auth 缺少有效的*gorm.DB对象")
	}
	return auth.gormDB
}

// NewContext 返回 ContextRepository
func (auth *Auth) NewContext(req *http.Request) *ContextRepository {
	return newContextRepository(req)
}

// RegisterProvider 注册登陆方式
func (auth *Auth) RegisterProvider(provider LoginProvider) {
	name := provider.Name()
	auth.Service.providers[name] = provider
	return
}

// LoginProvider 登陆方法
type LoginProvider interface {
	Name() string
	Login(credentials []byte) (user *User, err error)
}
