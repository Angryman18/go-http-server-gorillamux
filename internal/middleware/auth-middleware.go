package middleware

import (
	"context"
	"encoding/json"
	constants "go-server/internal/const"
	"go-server/internal/handler"
	utils "go-server/pkg/helper"
	"net/http"
	"time"
)

type Claim string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, err := utils.VerifyJWT(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errMsg := handler.NewError("JWT Token Error", err.Error())
			d, _ := json.Marshal(errMsg)
			w.Write(d)
			return
		}

		ctx := context.WithValue(r.Context(), constants.Claim, claim)
		timeOutCtx, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()
		r = r.WithContext(timeOutCtx)
		next.ServeHTTP(w, r)
	})
}
