package auth

import (
	"github.com/jinzhu/gorm"
)

// User 操作类...

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

// FindByOpenID 按openid查找
func (r *Repository) FindByOpenID(provider, openID string) (user *User, err error) {
	identity, err := r.FindIdentity(provider, openID)
	if err != nil {
		return nil, err
	}
	return r.Find(identity.UserID)
}

// FindByToken 根据token查找
func (r *Repository) FindByToken(tokenString string) (user *User, err error) {
	token, err := r.FindToken(tokenString)
	if err != nil {
		return nil, err
	}
	// load user
	return r.Find(token.UserID)
}

// Create 创建
func (r *Repository) Create(username, name string, avatar ...string) (u *User, err error) {
	user := User{
		Username: username,
		Name:     name,
	}
	if len(avatar) == 1 {
		user.Avatar = avatar[0]
	}
	err = r.db().Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
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
