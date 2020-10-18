package main

import (
	"fmt"

	"github.com/someshkoli/imageAPI/pkg/server"
)

func main() {
	srv := server.NewServer(8000)
	srv.Register()
	fmt.Println("Server listening at port 8000")
	if err := srv.Listen(); err != nil {
		fmt.Println(err)
	}
}
