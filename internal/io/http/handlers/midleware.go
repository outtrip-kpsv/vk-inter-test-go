package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"vk-inter-test-go/internal/io/ioutils"
	"vk-inter-test-go/internal/io/models"
	utilsJwt "vk-inter-test-go/internal/utils"
)

func (c *Controller) GlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.logger.Info("", zap.Reflect("req", r.URL))
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("authMiddleware", zap.Reflect("header", r.Header))
		if c.Bl.CheckJwt(r.Header.Get("Bearer")) {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			answer := models.ErrorResponse{
				Error: "wrong bearer",
			}
			ioutils.RespJson(w, answer)
			return
		}
	}
}

func (c *Controller) RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("CheckRole")
		bearer := r.Header.Get("Bearer")
		userName, err := utilsJwt.ExtractUsernameFromToken(bearer)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			answer := models.ErrorResponse{
				Error: "Token Extract Error :" + err.Error(),
			}
			ioutils.RespJson(w, answer)
			return
		}
		if c.Bl.CheckRole(userName, role) {
			next.ServeHTTP(w, r)

		} else {
			w.WriteHeader(http.StatusForbidden)
			answer := models.ErrorResponse{
				Error: "Access denied",
			}
			ioutils.RespJson(w, answer)
			return
		}
	}
}
