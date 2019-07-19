package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// respondJSON response with JSON format
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// deleteCookie 删除cookie
func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     name,
		Path:     "/",
		MaxAge:   -1,
	})
}

// setCookie 设置cookie
func setCookie(w http.ResponseWriter, name, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     name,
		Path:     "/",
		Value:    value,
		Expires:  expires,
	})
}
