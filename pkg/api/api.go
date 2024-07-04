package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/controllers"
	"github.com/kunalvirwal/go-mvc/pkg/middlewares"
	"github.com/kunalvirwal/go-mvc/pkg/models"
)

func StartRouter() {
	models.Connection()
	r := mux.NewRouter()
	CommonRoutes(r)

	r.Use(middlewares.AuthenticateToken) //////////complete thes 2 fn
	r.Use(middlewares.AuthorizeUser)

	log.Fatal(http.ListenAndServe(":4000", r))
}

func CommonRoutes(r *mux.Router) {

	r.HandleFunc("/", middlewares.CreateAdmin(controllers.LoginPage(false))).Methods("GET") ///////create Admin middleware
	r.HandleFunc("/signup", controllers.SignupPage(false)).Methods("GET")
	r.HandleFunc("/newUser", middlewares.SanitiseEmail(controllers.NewUser)).Methods("POST")
	r.HandleFunc("/login", middlewares.SanitiseEmail(controllers.Logging)).Methods("POST") ///////////////////////////////////add login

}

// func abc(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Request reached")
// }
