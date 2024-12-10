package handler

import "net/http"

type Health struct {
	Status  string
	Message string
}

func (a *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := Health{Status: "Ok", Message: "Server is Healthy"}
	writeResponse(w, http.StatusOK, &health)
}
