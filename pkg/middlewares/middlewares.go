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
	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"github.com/kunalvirwal/go-mvc/pkg/views"
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

				tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
					err := godotenv.Load()
					utils.CheckNilErr(err, "Unable to load .env in AuthenticateToken")
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
					utils.DeleteJWT(w)
					// fmt.Println("Token reset")    // delete that invalid token
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
				// fmt.Println(data.UUID)
				if data.Email != "" && data.Name != "" && data.UUID > 0 && (data.Role == "admin" || data.Role == "user") { /////can check expiration date condition too
					// fmt.Println(data)
					next.ServeHTTP(w, r)
				} else {
					http.SetCookie(w, &http.Cookie{
						Name:    "token",
						Value:   "",
						Expires: time.Unix(0, 0),
						Path:    "/",
					})

					http.Redirect(w, r, "/", http.StatusSeeOther)

				}

			} else {
				next.ServeHTTP(w, r)
			}

		} else {
			// fmt.Println("unprotected", r.URL.String())
			next.ServeHTTP(w, r)
		}

	})

}

func CreateAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")
		if !models.AdminExist(db) {
			user := models.USER{
				UUID:          1,
				NAME:          "admin",
				EMAIL:         "admin@sdslabs.com",
				PHN_NO:        9999999999,
				PASSWORD:      utils.SaltNhash("A"),
				ROLE:          "admin",
				ADMIN_REQUEST: nil,
			}
			result := db.Create(&user)
			utils.CheckNilErr(result.Error, "Unable to create Admin")
			fmt.Println("Admin account created")
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

		// http.Redirect(w, r, "/", http.StatusSeeOther)
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
