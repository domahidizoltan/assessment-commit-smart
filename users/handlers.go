package users

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

func (UserHandler) ListUsers(ctx echo.Context) error {
	fmt.Println("list users")
	return nil
}

func (UserHandler) CreateUser(ctx echo.Context) error {
	fmt.Println("create user")
	return nil
}

func (UserHandler) DeleteUser(ctx echo.Context, id int) error {
	fmt.Println("delete user", id)
	return nil
}

func (UserHandler) GetUser(ctx echo.Context, id int) error {
	fmt.Println("get user", id)
	return nil
}

func (UserHandler) UpdateUser(ctx echo.Context, id int) error {
	fmt.Println("update user", id)
	return nil
}
