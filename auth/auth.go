package auth

import "github.com/jinzhu/gorm"

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

// NewRepository 返回 Repository
func (a *Auth) NewRepository() *Repository {
	return newRepository(a.gormDB)
}
