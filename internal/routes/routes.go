package routes

import (
	"github.com/gorilla/mux"
)

const (
	HOME         = "/"
	LOGIN        = "/login"
	SIGNUP       = "/signup"
	HEALTH       = "/health"
	CREATE_TODO  = "/create_todo"
	GET_ALL_TODO = "/get-all-todo"
	UPDATE_TODO  = "/update-todo"
	DELETE_TODO  = "/delete-todo"
)

type Router struct {
	Router *mux.Router
}
