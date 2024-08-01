package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kunalvirwal/go-mvc/pkg/middlewares"
	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
)

func StartRouter() {
	r := mux.NewRouter()

	CommonRoutes(r)
	UserRoutes(r)
	AdminRoutes(r)

	path, err := filepath.Abs(".")

	if err == nil {
		s := http.StripPrefix("/styles/", http.FileServer(http.Dir(path+"/pkg/templates/styles/")))
		r.PathPrefix("/styles/").Handler(s)
	}

	models.CheckForAdmin()
	r.Use(middlewares.AuthenticateToken)

	log.Fatal(http.ListenAndServe(getPort(), r))
}

func getPort() string {
	err := godotenv.Load()
	utils.CheckNilErr(err, "Unable to load .env")
	port := os.Getenv("PORT")
	return port
}
