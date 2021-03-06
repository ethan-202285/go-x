package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/goodwong/go-x/auth/templates"
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
	userID := NewContext(r.Context()).UserID()
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		//case "GET":
		//// 出于安全考虑，默认不启用登陆演示页面
		//// 若需要启用，请自行添加，如：
		//// r.Get("/api/login", auths.Handler.LoginDemoPage)
		//	h.LoginDemoPage(w, r)

		case "PUT":
			h.HandleRenew(w, r)

		case "POST":
			h.HandleLogin(w, r)

		case "DELETE":
			// 这里ParseToken 可能会和 application层面的ParseToken重复，但是关系不大
			// 至少保证不能漏
			h.auth.Middleware.ParseToken(http.HandlerFunc(h.HandleLogout)).ServeHTTP(w, r)

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}

// LoginDemoPage 登陆、续约、注销的示例页面
func (h *Handler) LoginDemoPage(w http.ResponseWriter, r *http.Request) {
	// 查询数据
	// *应用开发不应该使用这个接口
	// *而应该专门设计返回用户信息的接口，如用GraphQL
	if r.URL.Query().Get("mode") == "data" {
		user := NewContext(r.Context()).User()

		// 已登陆
		if user != nil {
			respondJSON(w, map[string]interface{}{
				"id":   user.ID,
				"name": user.Name,
			}, http.StatusOK)
			return
		}

		// 未登录
		respondJSON(w, nil, http.StatusOK)
		return
	}

	// 显示页面
	// tpl := template.Must(template.ParseFiles("template/password_login.html"))
	html := templates.MustAsset("password_login.html")
	tpl := template.Must(template.New("password_login.html").Parse(string(html)))
	variables := map[string]string{}
	tpl.Execute(w, variables)
}
