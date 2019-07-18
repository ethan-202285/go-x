package auth

import (
	"encoding/json"
	"errors"
)

// UserIndentity 操作类...

// FindIdentity 查找UserIdentity
func (r *Repository) FindIdentity(provider, openID string) (indentity *UserIdentity, err error) {
	indentity = &UserIdentity{}
	err = r.db().Where("provider = ? and open_id = ?", provider, openID).Take(indentity).Error
	if err != nil {
		return nil, err
	}
	return
}

// FindIdentityByUser 查找UserIdentity
func (r *Repository) FindIdentityByUser(userID uint64, provider string) (indentity *UserIdentity, err error) {
	indentity = &UserIdentity{}
	err = r.db().Where("user_id = ? and provider = ?", userID, provider).Take(indentity).Error
	if err != nil {
		return nil, err
	}
	return
}

// CreateIdentity 创建UserIdentity
func (r *Repository) CreateIdentity(userID uint64, provider, openID string, data ...interface{}) (indentity *UserIdentity, err error) {
	if userID == 0 {
		return nil, errors.New("userID不能为0")
	}
	indentity = &UserIdentity{
		UserID:   userID,
		Provider: provider,
		OpenID:   openID,
	}
	if len(data) == 1 {
		bytes, _ := json.Marshal(data[0])
		jsonData := json.RawMessage(bytes)
		indentity.Data = &jsonData
	}
	err = r.db().Create(indentity).Error
	if err != nil {
		return nil, err
	}
	return
}

// UpdateIdentityData 更新UserIdentity
func (r *Repository) UpdateIdentityData(identity *UserIdentity, data interface{}) { // 只能更新Data数据
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = r.db().Model(identity).Update("data", json.RawMessage(bytes)).Error
	if err != nil {
		panic(err)
	}
}

// UpdateIdentityUser 更新UserIdentity
func (r *Repository) UpdateIdentityUser(identity *UserIdentity, user *User) { // 更新绑定的用户，单独出来接口，避免误操作
	err := r.db().Model(identity).Update("user_id", user.ID).Error
	if err != nil {
		panic(err)
	}
}
