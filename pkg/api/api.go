package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kunalvirwal/go-mvc/pkg/middlewares"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
)

func StartRouter() {
	r := mux.NewRouter()

	CommonRoutes(r)
	UserRoutes(r)
	AdminRoutes(r)

	r.Use(middlewares.AuthenticateToken)
	r.Use(middlewares.AuthorizeUser)

	log.Fatal(http.ListenAndServe(getPort(), r))
}

func getPort() string {
	err := godotenv.Load()
	utils.CheckNilErr(err, "Unable to load .env")
	port := os.Getenv("PORT")
	return port
}
