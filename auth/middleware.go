package auth

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

// newMiddleware 创建
func newMiddleware(auth *Auth) *Middleware {
	return &Middleware{auth: auth}
}

// Middleware 中间件
type Middleware struct {
	auth *Auth
}

// ParseToken 解析Token
// 解析token，在r.Context()里带上userID
// 如果没有登录，也不会强制要求登陆
// 如果需要强制要求登陆，
// 需要后面再加`Authenticated` 或 `AuthenticatedWithUser`
func (m *Middleware) ParseToken(next http.Handler) http.Handler {
	check := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			// jwt_token有效，直接过！
			if err == nil && token != nil && token.Valid && m.auth.Service.JwtInvalid(token) == false {
				// 带上userID继续
				claims := token.Claims.(jwt.MapClaims)
				userID := uint64(claims["sub"].(float64))
				ctx := NewContext(r.Context()).WithUserID(userID)
				next.ServeHTTP(w, ctx.AttachRequest(r))
				return
			}

			// jwt token无效，
			// 此时试图用refresh_token来续约jwt
			// 如果refresh_token 无效！让用户重新登录
			cookie, _ := r.Cookie("refresh_token")
			if cookie == nil || cookie.Value == "" {
				// 无效的refresh_token
				deleteCookie(w, "jwt")
				deleteCookie(w, "refresh_token")

				next.ServeHTTP(w, r)
				return
			}

			// 续约
			user, tokens, err := m.auth.Service.Renew(cookie.Value)
			if err != nil {
				// 续约失败
				deleteCookie(w, "jwt")
				deleteCookie(w, "refresh_token")

				next.ServeHTTP(w, r)
				return
			}

			// 成功续约！
			// 设置cookie
			setCookie(w, "jwt", tokens.Token, tokens.TokenExpires)

			// 带上userID继续
			ctx := NewContext(r.Context()).WithUser(user)
			next.ServeHTTP(w, ctx.AttachRequest(r))
		})
	}
	return jwtauth.Verifier(m.auth.jwtauth)(check(next))
}

// Authenticated 验证已登陆
func (m *Middleware) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(r.Context())
		if ctx.UserID() == 0 {
			respondJSON(
				w,
				map[string]string{"error": http.StatusText(http.StatusUnauthorized)},
				http.StatusUnauthorized,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthenticatedWithUser 验证已登录且自动附上User对象
func (m *Middleware) AuthenticatedWithUser(next http.Handler) http.Handler {
	attach := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(r.Context())

			user := ctx.User()
			// 如果前面已经通过续约登陆了，已经有用户数据了，
			// 这里就不需要重复查询数据库了
			if user != nil && user.Username != "" {
				next.ServeHTTP(w, r)
				return
			}

			// load
			user, err := m.auth.Repository.Find(ctx.UserID())
			// 这种情况一般是数据库直接删除用户导致的，
			// token验证通过，但是数据库找不到人
			// 按照失败处理：
			// 清理cookie
			if err == gorm.ErrRecordNotFound {
				deleteCookie(w, "jwt")
				deleteCookie(w, "refresh_token")
			}
			if err != nil {
				respondJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}

			// attach
			ctx.WithUser(user)
			next.ServeHTTP(w, ctx.AttachRequest(r))
		})
	}
	return m.Authenticated(attach(next))
}

// Authorized 必须是XX角色之一
//func (m *Middleware) Authorized(roles ...string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return func(next http.Handler) http.Handler {
//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				if role in roles {
//					...
//				}
//				next.ServeHTTP(w, r)
//			})
//		}(next)
//	}
//}
