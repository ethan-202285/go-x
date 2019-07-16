package auth

import (
	"time"
)

// Token操作类...

// FindToken 查找Token
func (r *Repository) FindToken(tokenString string) (token *Token, err error) {
	// parse and load token
	token, err = parseTokenString(tokenString, r.auth.secretKey)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if r.db().Where(Token{ID: token.ID}).Take(&token).RecordNotFound() {
		return nil, ErrInvalidToken
	}

	// verify token
	if !token.Verify() {
		return nil, ErrInvalidToken
	}
	if token.ExpiredAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	// valid
	return token, nil
}

// CreateToken 创建Token和tokenString
func (r *Repository) CreateToken(
	userID uint64, device string, remark ...string,
) (
	token *Token, tokenString string, err error,
) {
	// new
	params := Token{
		UserID: userID,
		Device: device,
	}
	if len(remark) == 1 {
		params.Remark = remark[0]
	}
	token = newToken(params, r.auth.secretKey)

	// create
	err = r.db().Create(token).Error
	if err != nil {
		return nil, "", err
	}

	return token, token.Stringify(), nil
}

// DeleteToken 删除Token
func (r *Repository) DeleteToken(userID uint64, device string) (err error) {
	return r.db().Where("user_id = ? and device = ?", userID, device).Delete(&Token{}).Error
}
