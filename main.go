package main

import (
	"commitsmart/users"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	var userHandler users.UserHandler
	e := echo.New()
	users.RegisterHandlersWithBaseURL(e, &userHandler, "/api/v1")

	fmt.Println("starting server on port 8000...")
	if err := e.Start("0.0.0.0:8000"); err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
