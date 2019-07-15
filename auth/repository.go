package auth

import (
	"github.com/jinzhu/gorm"
)

// newRepository 类
func newRepository(db *gorm.DB) *Repository {
	return &Repository{
		gormDB: db,
	}
}

// Repository 类
type Repository struct {
	gormDB *gorm.DB
}

func (r *Repository) db() *gorm.DB {
	if r.gormDB == nil {
		panic("Repository缺少有效的*gorm.DB对象")
	}
	return r.gormDB
}

// Find 根据id查找
func (r *Repository) Find(id uint64) (user *User, err error) {
	user = &User{ID: id}
	err = r.db().Take(user).Error
	if err != nil {
		return nil, err
	}
	return
}

// FindByUsername 根据用户名查找
func (r *Repository) FindByUsername(username string) (user *User, err error) {
	user = &User{}
	err = r.db().Where("username = ?", username).Take(user).Error
	if err != nil {
		return nil, err
	}
	return
}

// FindByToken 根据token查找
// TODO 待实现
func (r *Repository) FindByToken(token string) (user *User, err error) {
	panic("auth.FindByToken()功能未实现")
}

// Create 创建
func (r *Repository) Create(u *User) (err error) {
	return r.db().Create(u).Error
}

// Update 更新并保存
// (限制只能更新Name和Avatar)
func (r *Repository) Update(u *User, update User) error {
	update = User{
		Name:   update.Name,
		Avatar: update.Avatar,
	}
	return r.db().Model(u).Updates(update).Error
}

// UpdateUserName 更新用户名（独立出来，避免误用）
func (r *Repository) UpdateUserName(u *User, username string) error {
	update := User{
		Username: username,
	}
	return r.db().Model(u).Updates(update).Error
}

// List 查询列表
func (r *Repository) List(offset, limit int, where ...interface{}) (users []*User, err error) {
	var db *gorm.DB
	switch len(where) {
	case 0:
		db = r.db()
	case 1:
		db = r.db().Where(where[0])
	default:
		db = r.db().Where(where[0], where[1:]...)
	}
	err = db.Offset(offset).Limit(limit).Find(&users).Error
	return
}

// Count 统计文件长度
func (r *Repository) Count(where ...interface{}) (count int, err error) {
	err = r.db().Model(&User{}).Count(&count).Error
	return
}

// AutoMigrate 创建数据表
func (r *Repository) AutoMigrate() error {
	return r.db().AutoMigrate(&User{}, &UserIdentity{}, &Token{}, &UserLog{}).Error
}
