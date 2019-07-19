package auth

import (
	"net/http"

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

// Authenticated 验证已登陆
func (m *Middleware) Authenticated(next http.Handler) http.Handler {
	check := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			// jwt_token有效，直接过！
			if err == nil && token != nil && token.Valid {
				// 带上userID继续
				claims := token.Claims.(jwt.MapClaims)
				userID := uint64(claims["sub"].(float64))
				context := m.auth.NewContext(r).WithUserID(userID)
				next.ServeHTTP(w, context.Request())
				return
			}

			// jwt token无效，
			// 此时试图用refresh_token来续约jwt
			// 如果refresh_token 无效！让用户重新登录
			cookie, _ := r.Cookie("refresh_token")
			if cookie == nil || cookie.Value == "" {
				respondJSON(
					w,
					map[string]string{"error": http.StatusText(http.StatusUnauthorized)},
					http.StatusUnauthorized,
				)
				return
			}

			// 续约
			user, tokens, err := m.auth.Service.Renew(cookie.Value)
			if err != nil {
				deleteCookie(w, "refresh_token")
				respondJSON(w, map[string]string{"error": err.Error()}, http.StatusUnauthorized)
				return
			}

			// 设置cookie
			setCookie(w, "jwt", tokens.Token, tokens.TokenExpires)

			// 带上userID继续
			context := m.auth.NewContext(r).WithUser(user)
			next.ServeHTTP(w, context.Request())
		})
	}
	return jwtauth.Verifier(m.auth.jwtauth)(check(next))
}

// AuthenticatedWithUser 已登录且自动附上User对象
func (m *Middleware) AuthenticatedWithUser(next http.Handler) http.Handler {
	attach := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context := m.auth.NewContext(r)

			// load
			user, err := m.auth.Repository.Find(context.UserID())
			if err != nil {
				respondJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}

			// attach
			context.WithUser(user)
			next.ServeHTTP(w, context.Request())
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
