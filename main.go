package main

import (
	"context"
	"fmt"
	"go-server/internal/db"
	handler "go-server/internal/handler"
	"go-server/internal/middleware"
	"go-server/internal/routes"
	utils "go-server/pkg/helper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go-server/docs"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Todo Swagger
//	@version		1.0
//	@description	Todo Application by Shyam Mahanta.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		petstore.swagger.io
// @BasePath	/v2
func main() {

	route := mux.NewRouter()
	utils.LoadEnv()

	port := os.Getenv("PORT")

	if stringUtils.IsBlank(port) {
		log.Fatal("Port cannot be blank")
		os.Exit(1)
	}

	conn := db.Connect()

	apiV1 := "/api/v1"

	route.Use(middleware.LoggerMiddleware)

	route.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	handler := handler.AuthHandler{Conn: conn}
	publicRoutes := route.NewRoute().PathPrefix(apiV1).Subrouter()
	publicRoutes.Handle(routes.LOGIN, http.HandlerFunc(handler.Login)).Methods(http.MethodPost)
	publicRoutes.Handle(routes.SIGNUP, http.HandlerFunc(handler.SignupHandler)).Methods(http.MethodPost)
	publicRoutes.Handle(routes.HEALTH, http.HandlerFunc(handler.Health)).Methods(http.MethodGet)

	privateRoutes := route.NewRoute().PathPrefix(apiV1).Subrouter()
	privateRoutes.Use(middleware.AuthMiddleware)
	privateRoutes.Handle(routes.CREATE_TODO, http.HandlerFunc(handler.CreateTodo)).Methods(http.MethodPost)
	privateRoutes.Handle(routes.GET_ALL_TODO, http.HandlerFunc(handler.GetAllTodo)).Methods(http.MethodGet)
	privateRoutes.Handle(routes.UPDATE_TODO, http.HandlerFunc(handler.UpdateTodo)).Methods(http.MethodPost)
	privateRoutes.Handle(routes.DELETE_TODO, http.HandlerFunc(handler.DeleteTodo)).Methods(http.MethodPost)

	srvTxt := fmt.Sprintf("127.0.0.1:%s", port)
	srv := &http.Server{
		Handler:      route,
		Addr:         srvTxt,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server is running at ", srvTxt)
		fmt.Println(srv.ListenAndServe())
	}()

	<-sig

	ctx := context.Background()
	defer db.Close(conn)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Error Shutting Down Server")
		return
	}

	fmt.Println("Server Gracefully Shutdown")
}
