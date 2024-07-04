package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"github.com/kunalvirwal/go-mvc/pkg/views"
)

var user types.ContextKeyType = "abc"

func LoginPage(tried bool) http.HandlerFunc { ////////////////also check for pre
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, ok := r.Context().Value(user).(*types.Claims)
		if ok && data != nil {
			if data.Role == "admin" {
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
				return
			}
			if data.Role == "user" {
				http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
				return
			}
		}
		// fmt.Println("rendering login page")
		t := views.LoginView()
		t.Execute(w, tried)
	})
}

func SignupPage(exists bool) http.HandlerFunc { ////////////////also check for pre
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t := views.SignupView()
		t.Execute(w, exists)
	})
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := strings.TrimSpace(r.FormValue("name"))
	password := strings.TrimSpace(r.FormValue("password"))
	email := strings.TrimSpace(r.FormValue("email"))
	phn_no, err := strconv.ParseInt(r.FormValue("phn_no"), 10, 64)
	if err != nil || !utils.Sanitise(name, false) || phn_no > 9999999999 || phn_no < 1000000000 || len(name) == 0 || len(name) > 50 {
		t := views.SignupView()
		t.Execute(w, true)
		return
	}
	// var user models.USER
	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance")
	_, found := models.SearchUserEmail(db, email)
	if found {
		t := views.SignupView()
		t.Execute(w, true)
		return
	}
	models.CreateNewUser(db, name, email, phn_no, password) ////////////////////saltnhash
	Logging(w, r)
}

func Logging(w http.ResponseWriter, r *http.Request) {
	password := strings.TrimSpace(r.FormValue("password"))
	email := strings.TrimSpace(r.FormValue("email"))
	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance")
	user, found := models.SearchUserEmail(db, email)

	if found && user.PASSWORD == password { /////////////////////////// compare salted passwords

		token, expirationTime := utils.GenerateJWT(user.UUID, email, user.NAME, user.ROLE)

		fmt.Println(token)

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})

		fmt.Println("JWT set")

	} else {
		t := views.LoginView()
		t.Execute(w, true)
	}
}
