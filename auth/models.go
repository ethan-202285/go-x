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
	Provider string           `gorm:"primary_key"`
	OpenID   string           `gorm:"primary_key"` // 第三方主键
	Data     *json.RawMessage ``                   // json原始数据（密码也可以放这里）
}

// UserLogout 主动注销行为（数据量较小，可内存保存）
type UserLogout struct {
	UserID       uint64
	LocationCode string
	LogoutAt     time.Time //（只要是userID在这里的，并且user.iat < LogoutAt的，都要重新登录）
}

// TableName 按照gorm方式制定数据库表名
// ** 这里需要阻止建表
func (u *UserLogout) TableName() string {
	panic("UserLogout是内存表，无需创建数据库")
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
