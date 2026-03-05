package handler

import "net/http"

func SetupRoutes(authHandler *AuthHandler, jobHandler *JobHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/jobs", jobHandler.CreateJob)
	mux.HandleFunc("/jobs/status", jobHandler.GetJobStatus)

	return mux
}
