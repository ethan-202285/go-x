package auth

import (
	"context"
	"net/http"
)

func newContextRepository(req *http.Request) *ContextRepository {
	return &ContextRepository{request: req}
}

// ContextRepository 类
type ContextRepository struct {
	request *http.Request
}

// UserID 在context里获取UserID
func (r *ContextRepository) UserID() uint64 {
	user := r.User()
	if user == nil {
		return 0
	}
	return user.ID
}

// User 在context里获取User
func (r *ContextRepository) User() *User {
	value := r.request.Context().Value(contextKeyUser)
	user, ok := value.(*User)
	if !ok {
		return nil
	}
	return user
}

// WithUserID 在context里带上UserID
func (r *ContextRepository) WithUserID(userID uint64) *ContextRepository {
	// 检查是否已经有User
	user := r.User()
	if user == nil || user.ID != userID {
		user := &User{ID: userID}
		ctx := r.request.Context()
		ctx = context.WithValue(ctx, contextKeyUser, user)
		r.request = r.request.WithContext(ctx)
	}
	return r
}

// WithUser 在context里带上User
func (r *ContextRepository) WithUser(user *User) *ContextRepository {
	ctx := r.request.Context()
	ctx = context.WithValue(ctx, contextKeyUser, user)
	r.request = r.request.WithContext(ctx)
	return r
}

// Request 返回 Request
func (r *ContextRepository) Request() *http.Request {
	return r.request
}

var (
	contextKeyUser = &contextKey{"user"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "auth context key " + k.name
}
