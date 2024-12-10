package handler

import (
	"context"
	"go-server/internal/validation"
	utils "go-server/pkg/helper"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Conn *pgxpool.Pool
}

type SignupResponseData struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Id       string `json:"id"`
}

type SignUpResponse struct {
	Data *SignupResponseData `json:"data"`
}

// @title Login API
// @Consume json
// @Tags Authentication
// @Success 200 {object} SignUpResponse
// @Produce json
// @Param request body validation.LoginReqestBody true "Login Credentials"
// @Router /api/v1/login [post]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	payload, err := validation.ValidateLoginInfo(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Invalid Payload", err.Error()))
		return
	}
	var UserInfo UserInfo
	queryErr := a.Conn.QueryRow(r.Context(), `select id, name, password from users where email = $1`, payload.Email).Scan(&UserInfo.Id, &UserInfo.Name, &UserInfo.Password)
	if queryErr != nil {
		writeResponse(w, http.StatusNotFound, NewError("User not found", queryErr.Error()))
		return
	}
	if hashErr := bcrypt.CompareHashAndPassword([]byte(UserInfo.Password), []byte(payload.Password)); hashErr != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Password is not correct", hashErr.Error()))
		return
	}
	token, _ := utils.CreateJWT(UserInfo.Id)
	UserResp := SignUpResponse{Data: &SignupResponseData{Token: token, Name: UserInfo.Name, Email: payload.Email, Id: UserInfo.Id}}
	writeResponse(w, http.StatusOK, UserResp)
}

// @title Signup API
// @Consume json
// @Tags Authentication
// @Success 200 {object} SignupResponseData
// @Produce json
// @Param request body validation.RequestBody true "Signup Credentials"
// @Router /api/v1/signup [post]
func (u *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := validation.ValidateSignupInfo(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Invalid payload", err.Error()))
		return
	}

	genPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Error Creating Password", err.Error()))
		return
	}
	var resp SignupResponseData
	userId := uuid.New().String()
	sqError := u.Conn.QueryRow(context.Background(), `insert into users (id, name, email, password) values ($1, $2, $3, $4) returning id, name, email`, userId, body.Name, body.Email, genPassword).Scan(&resp.Id, &resp.Name, &resp.Email)
	if sqError != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Failed to signup", sqError.Error()))
		return
	}
	str, err := utils.CreateJWT(resp.Id)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, NewError("Unauthorized", err.Error()))
		return
	}
	resp.Token = str
	writeResponse(w, http.StatusOK, resp)

}
