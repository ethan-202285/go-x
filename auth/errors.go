package auth

import (
	"errors"
)

// ErrSaveConflict 冲突错误
var ErrSaveConflict = errors.New("更新失败，数据已过期")
