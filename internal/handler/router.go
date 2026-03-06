package handler

import (
	"net/http"

	"github.com/alialin/scraperq/internal/middleware"
)

func SetupRoutes(authHandler *AuthHandler, jobHandler *JobHandler, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/jobs", middleware.AuthMiddleware(jwtSecret, jobHandler.CreateJob))
	mux.HandleFunc("/jobs/status", middleware.AuthMiddleware(jwtSecret, jobHandler.GetJobStatus))

	return mux
}
