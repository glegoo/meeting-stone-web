package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(cfg Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Dev 直连（VITE_API_ORIGIN）时需要 CORS；生产同域部署一般不需要。
	if cfg.AppEnv != "prod" || cfg.WebOrigin != "" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   allowedOrigins(cfg),
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		Ok(w, map[string]any{
			"ok":      true,
			"version": "v0",
		})
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(WithAuth)

		api.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			Ok(w, map[string]any{
				"ok": true,
			})
		})

		api.Get("/games", handleGetGames)
		api.Get("/games/{gameKey}/manifest", handleGetManifest)
		api.Get("/games/{gameKey}/activities", handleGetActivities)
		api.Get("/games/{gameKey}/lobby/activities", handleGetLobbyActivities)
	})

	return r
}

func allowedOrigins(cfg Config) []string {
	if cfg.WebOrigin != "" {
		return []string{cfg.WebOrigin}
	}
	// 开发默认只放行本地地址，避免不必要的全开放。
	return []string{"http://localhost:*", "http://127.0.0.1:*"}
}
