package io

import (
	"net/http"
	"vk-inter-test-go/internal/io/http/handlers"
	"vk-inter-test-go/internal/io/ioutils"
)

func SetupRoutes(contr *handlers.Controller) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/create/user", contr.CreateUser)
	mux.HandleFunc("/api/login", contr.AuthUser)

	mux.HandleFunc("/api/actor", contr.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			contr.GetAllActors(w, r)
		case http.MethodPost:
			contr.RequireRole("admin", contr.CreateActor)(w, r)
		case http.MethodDelete:
			contr.RequireRole("admin", contr.DeleteActor)(w, r)
		case http.MethodPatch:
			contr.RequireRole("admin", contr.UpdateActor)(w, r)

		default:
			ioutils.HandleInvalidMethodResponse(w, r.Method)
		}
	}))
	mux.HandleFunc("/api/movie", contr.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			contr.GetMovies(w, r)
		case http.MethodPost:
			contr.RequireRole("admin", contr.CreateMovie)(w, r)
		case http.MethodDelete:
			contr.RequireRole("admin", contr.DeleteMovie)(w, r)
		case http.MethodPatch:
			contr.RequireRole("admin", contr.UpdateMovie)(w, r)

		default:
			ioutils.HandleInvalidMethodResponse(w, r.Method)
		}
	}))

	muxN := use(mux, contr.GlobalMiddleware)

	return muxN
}

func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}
