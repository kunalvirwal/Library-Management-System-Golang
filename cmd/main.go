package main

import (
	"fmt"

	"github.com/kunalvirwal/go-mvc/pkg/api"
)

func main() {
	fmt.Println("Server started!")
	api.StartRouter()
	fmt.Println("Hello")
}
