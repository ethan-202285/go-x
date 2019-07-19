package auth

import (
	"errors"
)

// ErrInvalidToken 无效的 RefreshToken
var ErrInvalidToken = errors.New("无效的 RefreshToken")
