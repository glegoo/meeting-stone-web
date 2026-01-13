package httpserver

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

// WithAuth 是一个最小占位鉴权：解析 Bearer token，但暂不校验。
// 后续接入 JWT 校验/refresh cookie 流程时替换这里即可。
func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
			if token != "" {
				// 先把 token 原样塞进 context，当作 user_id 占位。
				ctx := context.WithValue(r.Context(), userIDKey, token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) (string, bool) {
	v := r.Context().Value(userIDKey)
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}
