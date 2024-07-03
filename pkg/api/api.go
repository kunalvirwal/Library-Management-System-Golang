package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/controllers"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
)

func StartRouter() {
	r := mux.NewRouter()

	CommonRoutes(r)

	//common middlewares
	r.Use(utils.AuthenticateToken) //////////complete thes 2 fn
	r.Use(utils.AuthorizeUser)

	log.Fatal(http.ListenAndServe(":4000", r))
}

func CommonRoutes(r *mux.Router) {

	r.HandleFunc("/", controllers.LoginPage(false)).Methods("GET") ///////create Admin middleware
	r.HandleFunc("/signup", controllers.SignupPage(false)).Methods("GET")
	r.HandleFunc("/login", utils.SanitiseEmail(abc)).Methods("POST") // demo

}

func abc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request reached")
}
