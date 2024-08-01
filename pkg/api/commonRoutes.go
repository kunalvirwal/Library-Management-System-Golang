package api

import (
	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/controllers"
	"github.com/kunalvirwal/go-mvc/pkg/middlewares"
)

func CommonRoutes(r *mux.Router) {

	r.HandleFunc("/", controllers.LoginPage(false)).Methods("GET")
	r.HandleFunc("/signup", controllers.SignupPage(false)).Methods("GET")
	r.HandleFunc("/newUser", middlewares.SanitiseEmail(controllers.NewUser)).Methods("POST")
	r.HandleFunc("/login", middlewares.SanitiseEmail(controllers.Logging)).Methods("POST")
	r.HandleFunc("/books", controllers.GetBookCatalog).Methods("GET")
	r.HandleFunc("/books/{buid}", controllers.GetBookPage).Methods("GET")
	r.HandleFunc("/checkout", controllers.MakeCheckoutReq).Methods("POST")
	r.HandleFunc("/checkin", controllers.MakeCheckinReq).Methods("POST")
	r.HandleFunc("/pending", controllers.GetPending).Methods("GET")
	r.HandleFunc("/account", controllers.Account).Methods("GET")
	r.HandleFunc("/account", controllers.EditAccount).Methods("POST")
	r.HandleFunc("/cvt_admin", controllers.GetCvtAdmin).Methods("GET")
	r.HandleFunc("/cvt_admin", controllers.PostCvtAdmin).Methods("POST")

	r.HandleFunc("/logout", controllers.Logout)

}
