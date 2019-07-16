package auth

import (
	"errors"
)

// ErrSaveConflict 冲突错误
var ErrSaveConflict = errors.New("更新失败，数据已过期")

// ErrInvalidToken 无效的 RememberToken
var ErrInvalidToken = errors.New("无效的 RememberToken")
