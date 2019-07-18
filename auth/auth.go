package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
)

// New 返回Auth类
func New(config Config) *Auth {
	auth := &Auth{
		secretKey: config.SecretKey,
		gormDB:    config.DB,
		jwtauth:   jwtauth.New("HS256", config.SecretKey, nil),
	}
	auth.Repository = newRepository(auth)
	auth.Service = newService(auth)
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
