package auth

import (
	"context"
	"net/http"
)

// NewContext 新建context装饰类
func NewContext(ctx context.Context) *ContextRepository {
	return &ContextRepository{context: ctx}
}

// ContextRepository 类
type ContextRepository struct {
	context context.Context
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
	value := r.context.Value(contextKeyUser)
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
		r.context = context.WithValue(r.context, contextKeyUser, user)
	}
	return r
}

// WithUser 在context里带上User
func (r *ContextRepository) WithUser(user *User) *ContextRepository {
	r.context = context.WithValue(r.context, contextKeyUser, user)
	return r
}

// AttachRequest 返回 Request
func (r *ContextRepository) AttachRequest(req *http.Request) *http.Request {
	return req.WithContext(r.context)
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
