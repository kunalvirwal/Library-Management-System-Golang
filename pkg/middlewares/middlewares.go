package middlewares

import (
	"context"
	"fmt"
	"os"
	"time"

	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"github.com/kunalvirwal/go-mvc/pkg/views"
)

var user types.ContextKeyType = "abc"

// decodes the token if present and path is protected and thus verifies JWT
func AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // route specific middlewares have to return http.HandlerFunc

		// if not an unprotected path and cookie present then verify jwt
		var unprotectedPaths = [3]string{"/signup", "/newUser", "/login"}
		var protected bool = true

		// check if it is a protected path
		for _, path := range unprotectedPaths {
			if path == r.URL.String() {
				protected = false
			}
		}

		// checks for protected paths
		if protected {
			cookie, err := r.Cookie("token")

			// here we didnt use checkNillErr because that fun causes a panic and is only for server side errors
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			token := cookie.Value
			claims := &types.Claims{}
			if token != "" {
				tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
					err := godotenv.Load()
					utils.CheckNilErr(err, "Unable to load .env in AuthenticateToken")
					return []byte(os.Getenv("SECRET_KEY")), nil
				})
				if err != nil || !tkn.Valid {
					fmt.Println("Invalid JWT token", err)
					http.Redirect(w, r, "/", http.StatusBadRequest)
					fmt.Println("Token reset")

				}
				if claims.Email != "" && claims.Name != "" && claims.UUID > 0 && (claims.Role == "admin" || claims.Role == "user") {
					user = "abc"
					ctx := context.WithValue(r.Context(), user, claims)
					r = r.WithContext(ctx)
				} else {
					http.SetCookie(w, &http.Cookie{
						Name:    "token",
						Value:   "",
						Expires: time.Unix(0, 0),
						Path:    "/",
					})
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func SanitiseEmail(next http.HandlerFunc) http.HandlerFunc { // route specific middlewares have to return http.HandlerFunc

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			fmt.Println("Invalid Post req paramemeters")
			return
		}
		inpEmail := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		valid := true
		parts := strings.Split(inpEmail, "@")
		if len(parts) == 2 {
			valid = utils.Sanitise(parts[1], true) && !strings.Contains(parts[0], " ")
		} else {
			valid = false
		}
		if valid {
			next.ServeHTTP(w, r)
			return
		}
		t := views.LoginView()
		t.Execute(w, true)
	})
}

func IsUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, ok := r.Context().Value(user).(*types.Claims)
		if data != nil && ok {
			if data.Role == "user" {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, ok := r.Context().Value(user).(*types.Claims)
		if data != nil && ok {
			if data.Role == "admin" {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
