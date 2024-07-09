package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

func CheckNilErr(err error, st string) {

	if err != nil {
		fmt.Println(st)
		panic(err)
	}
}

func Sanitise(query string, space bool) bool {
	wrongCharacters := []string{"'", "\"", "`", "--", "=", "(", ")", ","}
	if space {
		wrongCharacters = append(wrongCharacters, " ")
	}
	for _, val := range wrongCharacters {
		if strings.Contains(query, val) {
			return false
		}
	}
	return true
}

func GenerateJWT(uuid int, email string, name string, role string) (string, time.Time) {

	err := godotenv.Load()
	CheckNilErr(err, "Unable to load .env")
	secret_key := []byte(os.Getenv("SECRET_KEY"))

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &types.Claims{
		UUID:  uuid,
		Email: email,
		Name:  name,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret_key)
	CheckNilErr(err, "Unable to sign JWT key")

	return tokenString, expirationTime

}

func DeleteJWT(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
}

func MatchHashtoPassword(originalHash string, password string) bool {
	inp_pwd := []byte(password)
	og_pwd := []byte(originalHash)
	err := bcrypt.CompareHashAndPassword(og_pwd, inp_pwd)
	return err == nil
}

func SaltNhash(pwd string) string {
	inp_pwd := []byte(pwd)
	hashedPwdBytes, err := bcrypt.GenerateFromPassword(inp_pwd, bcrypt.MinCost)
	CheckNilErr(err, "Unable to Hash password")
	return string(hashedPwdBytes)
}
