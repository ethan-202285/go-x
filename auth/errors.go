package auth

import (
	"errors"
)

// ErrInvalidToken 无效的 RememberToken
var ErrInvalidToken = errors.New("无效的 RememberToken")
