package auth

import (
	"encoding/json"
	"time"
)

// User 类型
type User struct {
	ID        uint64
	Username  string `gorm:"unique_index;not null"`
	Name      string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserIdentity 登录方式(小程序、钉钉……密码都可以）
type UserIdentity struct {
	UserID   uint64
	Provider string           `gorm:"primary_key;not null"`
	OpenID   string           `gorm:"primary_key;not null"` // 第三方主键
	Data     *json.RawMessage ``                            // json原始数据（密码也可以放这里）
}

// UserLog 用户行为记录（登录、注销）
type UserLog struct {
	ID        uint64
	UserID    uint64
	Action    string //（程序决定）
	Remark    string //（备注）
	IP        string
	UA        string
	CreatedAt time.Time
}
