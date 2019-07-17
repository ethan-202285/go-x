package auth

import (
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
)

// New 返回Auth类
func New(config Config) *Auth {
	return &Auth{
		secretKey: config.SecretKey,
		gormDB:    config.DB,
	}
}

// Config 配置
type Config struct {
	SecretKey []byte
	DB        *gorm.DB
}

// Auth 认证类
type Auth struct {
	secretKey []byte
	gormDB    *gorm.DB
}

func (auth *Auth) db() *gorm.DB {
	if auth.gormDB == nil {
		panic("Auth 缺少有效的*gorm.DB对象")
	}
	return auth.gormDB
}

// NewRepository 返回 Repository
func (auth *Auth) NewRepository() *Repository {
	return newRepository(auth)
}

// NewContext 返回 ContextRepository
func (auth *Auth) NewContext(req *http.Request) *ContextRepository {
	return newContextRepository(req)
}
