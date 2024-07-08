package api

import (
	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/controllers"
	"github.com/kunalvirwal/go-mvc/pkg/middlewares"
)

func AdminRoutes(r *mux.Router) {
	r.HandleFunc("/admin/dashboard", middlewares.IsAdmin(controllers.GetAdminData)).Methods("GET")
	r.HandleFunc("/admin/approve/{uuid}/{buid}", middlewares.IsAdmin(controllers.Approve)).Methods("GET")
	r.HandleFunc("/admin/deny/{uuid}/{buid}", middlewares.IsAdmin(controllers.Deny)).Methods("GET")
	r.HandleFunc("/admin/editbook/{buid}", middlewares.IsAdmin(controllers.Editbook)).Methods("GET")
	r.HandleFunc("/admin/editbook/{buid}", middlewares.IsAdmin(controllers.SaveEditbookChanges)).Methods("POST")
	r.HandleFunc("/admin/deletebook/{buid}", middlewares.IsAdmin(controllers.DeleteBook)).Methods("GET")
	r.HandleFunc("/admin/addbook", middlewares.IsAdmin(controllers.NewBook)).Methods("GET")
	r.HandleFunc("/admin/addbook", middlewares.IsAdmin(controllers.SaveNewBook)).Methods("POST")

}
