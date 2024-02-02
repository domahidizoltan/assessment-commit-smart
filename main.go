package main

import (
	"commitsmart/users"
	userServer "commitsmart/users/generated"
	"log"

	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	var userHandler users.UserHandler
	e := echo.New()
	userServer.RegisterHandlersWithBaseURL(e, &userHandler, "/api/v1")

	fmt.Println("starting server on port 8000...")
	if err := e.Start("0.0.0.0:8000"); err != nil {
		log.Default().Printf("error listening for server: %s\n", err)
	}
}
