package utils

import (
	"fmt"
	"strings"
)

func CheckNilErr(err error, st string) {

	if err != nil {
		fmt.Println(st)
		panic(err)
	}
}

func sanitise(query string) bool {
	wrongCharacters := []string{"'", "\"", "`", "--", "=", " ", "(", ")", ","}
	for _, val := range wrongCharacters {
		if strings.Contains(query, val) {
			return false
		}
	}
	return true
}
