package utils

import (
	"fmt"
	"strings"
)

// global variables
var Admin = "admin"
var User = "user"

// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
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
