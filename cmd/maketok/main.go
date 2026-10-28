package main

import (
	"fmt"

	authpkg "github.com/sidalikharef/goremote/internal/auth"
)

func main() {
	s := authpkg.New("/tmp/test.auth.json")
	tok, err := s.GenerateBootstrapToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(tok)
}
