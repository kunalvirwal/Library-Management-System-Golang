package controllers

import (
	"net/http"

	"github.com/kunalvirwal/go-mvc/pkg/types"
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
