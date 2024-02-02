package users

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	server "commitsmart/users/generated"

	"github.com/labstack/echo/v4"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserHandler struct {
}

func init() {
	if err := mgm.SetDefaultConfig(nil, "local", options.Client().ApplyURI("mongodb://root:pass@localhost:27018/")); err != nil {
		panic(err)
	}
}

func (UserHandler) ListUsers(ctx echo.Context) error {
	users := []User{}
	if err := mgm.Coll(&User{}).SimpleFind(&users, bson.M{}); err != nil {
		log.Default().Printf("failed to list users: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	dtos := []server.User{}
	for _, u := range users {
		dtos = append(dtos, modelToDto(u))
	}

	return ctx.JSON(http.StatusOK, dtos)
}

func (UserHandler) CreateUser(ctx echo.Context) error {
	dto := new(server.User)
	if err := ctx.Bind(dto); err != nil {
		return ctx.String(http.StatusBadRequest, "Bad request")
	}
	user := &User{
		Email:    dto.Email,
		Password: dto.Password,
		Username: dto.Username,
	}
	if err := mgm.Coll(user).Create(user); err != nil {
		log.Default().Printf("failed to create user: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, user.ID.Hex())
}

func (UserHandler) DeleteUser(ctx echo.Context, id string) error {
	id = sanitizeIdString(id)
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Printf("failed to parse id %s: %s", id, err)
	}

	//TODO no error but still not deletes
	_, err = mgm.Coll(&User{}).DeleteOne(context.TODO(), bson.M{"_id": _id.String()})
	if err != nil {
		log.Default().Printf("failed to delete user %s: %s", id, err.Error())
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (UserHandler) GetUser(ctx echo.Context, id string) error {
	id = sanitizeIdString(id)
	user := &User{}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Printf("failed to parse id %s: %s", id, err.Error())
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if err := mgm.Coll(&User{}).FindByID(_id, user); err != nil {
		log.Default().Printf("failed to get user %s: %s", id, err.Error())
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	dto := modelToDto(*user)

	return ctx.JSON(http.StatusOK, dto)
}

func (UserHandler) UpdateUser(ctx echo.Context, id string) error {
	//TODO
	fmt.Println("update user", id)
	return nil
}

func modelToDto(u User) server.User {
	id := u.ID.String()
	return server.User{
		Email:    u.Email,
		Id:       &id,
		Username: u.Username,
		Password: u.Password,
	}
}

func sanitizeIdString(id string) string {
	id = strings.TrimSpace(id)
	id = strings.Trim(id, "\"")
	return id
}
