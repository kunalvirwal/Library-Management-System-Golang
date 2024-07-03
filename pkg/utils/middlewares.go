package utils

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
)

var user types.ContextKeyType = "abc"

// decodes the token if present and path is protected
func AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var unprotectedPaths = [3]string{"/signup", "/newUser", "/login"} // if not an unprotected path and cookie present then verify jwt
		var protected bool = true

		for _, path := range unprotectedPaths {
			if path == r.URL.String() {
				protected = false // check if it is a protected path
			}
		}

		if protected {
			cookie, err := r.Cookie("token")
			if err != nil { // here we didnt use checkNillErr because that fun causes a panic and is only for server side errors
				next.ServeHTTP(w, r)
				return
			}

			token := cookie.Value
			claims := &types.Claims{}

			if token != "" {
				fmt.Println("Entered")
				tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
					err := godotenv.Load()
					CheckNilErr(err, "Unable to load .env in AuthenticateToken")

					return []byte(os.Getenv("SECRET_KEY")), nil
				})
				if err != nil {
					fmt.Println("Invalid JWT token", err)
					fmt.Println("Token reset") // delete that invalid token

				}
				// fmt.Printf("%T \n %v", claims, claims)
				user = "abc"
				ctx := context.WithValue(r.Context(), user, claims)
				r = r.WithContext(ctx)

				if !tkn.Valid {
					http.SetCookie(w, &http.Cookie{
						Name:    "token",
						Value:   "",
						Expires: time.Unix(0, 0),
						Path:    "/",
					})
					// fmt.Println("Token reset")    // delete that invalid token
				} else {
					ctx := context.WithValue(r.Context(), user, claims)
					r = r.WithContext(ctx)

				}

			}

		}
		next.ServeHTTP(w, r)
	})
}

func AuthorizeUser(next http.Handler) http.Handler { // common to all middlewares utilised by r.use can return http handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var unprotectedPaths = [3]string{"/signup", "/newUser", "/login"}
		var protected = true
		for _, path := range unprotectedPaths {
			if path == r.URL.String() {
				protected = false // check if it is a protected path
			}
		}

		if protected {
			// fmt.Println("protected", r.URL.String())

			data, ok := r.Context().Value(user).(*types.Claims)
			if data != nil && ok {
				if data.Email != "" && data.Name != "" && (data.Role == "admin" || data.Role == "user") { /////can check expiration date condition too
					fmt.Println(data)
					next.ServeHTTP(w, r)
				}

			}
			///////////////////////

			next.ServeHTTP(w, r)

		} else {
			// fmt.Println("unprotected", r.URL.String())
			next.ServeHTTP(w, r)
		}

	})

}

func SanitiseEmail(next http.HandlerFunc) http.HandlerFunc { // route specific middlewares have to return http.HandlerFunc

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		inpEmail := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		valid := true
		parts := strings.Split(inpEmail, "@")
		if len(parts) == 2 {
			valid = sanitise(parts[1])
		} else {
			valid = false
		}
		if valid {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
