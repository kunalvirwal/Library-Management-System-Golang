package api

import (
	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/controllers"
	"github.com/kunalvirwal/go-mvc/pkg/middlewares"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/user/dashboard", middlewares.IsUser(controllers.GetUserData)).Methods("GET")
	r.HandleFunc("/user/refresh", controllers.Refresh).Methods("POST")

}
