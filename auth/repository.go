package auth

import (
	"github.com/jinzhu/gorm"
)

// newRepository 类
func newRepository(auth *Auth) *Repository {
	return &Repository{auth: auth}
}

// Repository 类
type Repository struct {
	auth *Auth
}

func (r *Repository) db() *gorm.DB {
	if r.auth == nil {
		panic("缺少auth字段")
	}
	return r.auth.db()
}

// AutoMigrate 创建数据表
func (r *Repository) AutoMigrate() error {
	return r.db().AutoMigrate(&User{}, &UserIdentity{}, &Token{}, &UserLog{}).Error
}
