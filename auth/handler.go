package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
)

// newHandler 创建
func newHandler(auth *Auth) *Handler {
	return &Handler{auth: auth}
}

// Handler 即Controller
type Handler struct {
	auth *Auth
}

// HandleLogin 登陆
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// 参数
	provider := r.URL.Query().Get("provider")
	device := r.URL.Query().Get("device")
	remember := false
	switch rememberParam := r.URL.Query().Get("remember"); rememberParam {
	case "1", "true":
		remember = true
	case "", "0", "false", "null":
		remember = false
	default:
		remember = false
	}
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, map[string]string{"error": "request body读取错误"}, http.StatusBadRequest)
		return
	}

	// 登陆逻辑
	tokens, err := h.auth.Service.Login(provider, payload, remember, device)
	if err != nil {
		respondJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// 设置cookie
	setCookie(w, "jwt", tokens.Token, tokens.TokenExpires)
	if tokens.RefreshToken != nil {
		setCookie(w, "refresh_token", *tokens.RefreshToken, *tokens.RefreshTokenExpires)
	}

	// 返回
	respondJSON(w, tokens, http.StatusOK)
}

// HandleRenew JWT续约
func (h *Handler) HandleRenew(w http.ResponseWriter, r *http.Request) {
	// 参数
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, map[string]string{"error": "request body读取错误"}, http.StatusBadRequest)
		return
	}
	var params struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		respondJSON(w, map[string]string{"error": "request body读取错误"}, http.StatusBadRequest)
		return
	}

	// 续约
	_, tokens, err := h.auth.Service.Renew(params.RefreshToken)
	if err != nil {
		respondJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// 设置cookie
	setCookie(w, "jwt", tokens.Token, tokens.TokenExpires)

	// 返回
	respondJSON(w, tokens, http.StatusOK)
}

// HandleLogout 登出
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	userID := h.auth.NewContext(r).UserID()
	if userID == 0 {
		respondJSON(w, "无需登出!", http.StatusOK)
		return
	}

	// 清理数据库
	user := &User{ID: userID}
	device := r.URL.Query().Get("device")
	h.auth.Service.Logout(user, device)

	// 清理cookie
	deleteCookie(w, "jwt")
	deleteCookie(w, "refresh_token")

	// 返回
	respondJSON(w, "登出成功!", http.StatusOK)
	return
}

// Mux 返回多路复用器
func (h *Handler) Mux() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/login", func(r chi.Router) {
		r.Post("/", h.HandleLogin)

		authenticated := h.auth.Middleware.AuthenticatedWithUser
		r.With(authenticated).Delete("/", h.HandleLogout)
	})

	return r
}
