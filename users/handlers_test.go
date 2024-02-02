package users

import (
	userServer "commitsmart/users/generated"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kamva/mgm/v3"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// testing with the same DB, this will delete all the records
func cleanDb(t *testing.T) {
	_, err := mgm.Coll(&User{}).DeleteMany(context.TODO(), bson.M{})
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	cleanDb(t)

	dtoJson := `{"email":"test@test.com", "password":"pass","username":"test"}`

	handler := &UserHandler{}
	e := echo.New()

	//create
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(dtoJson)))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(createReq, rec)

	require.NoError(t, handler.CreateUser(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
	createdUserId := rec.Body.String()
	assert.NotEmpty(t, createdUserId)

	//get
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/:id", nil)
	getReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(getReq, rec)
	c.SetParamNames("id")
	c.SetParamValues(createdUserId)

	require.NoError(t, handler.GetUser(c, createdUserId))
	assert.Equal(t, http.StatusOK, rec.Code)
	userJson := rec.Body.String()
	var fetchedUser userServer.User
	err := json.Unmarshal([]byte(userJson), &fetchedUser)
	require.NoError(t, err)
	assert.Equal(t, *fetchedUser.Email, "test@test.com")
	assert.Equal(t, *fetchedUser.Username, "test")
	assert.Equal(t, *fetchedUser.Password, "pass")

	//delete
	delReq := httptest.NewRequest(http.MethodDelete, "/api/v1/users/:id", nil)
	delReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(delReq, rec)
	c.SetParamNames("id")
	c.SetParamValues(createdUserId)

	require.NoError(t, handler.DeleteUser(c, createdUserId))
	assert.Equal(t, http.StatusNoContent, rec.Code)

	//list
	listReq := httptest.NewRequest(http.MethodDelete, "/api/v1/users", nil)
	listReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(listReq, rec)

	require.NoError(t, handler.ListUsers(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	usersJson := rec.Body.String()
	var fetchedUsers []userServer.User
	err = json.Unmarshal([]byte(usersJson), &fetchedUsers)
	require.NoError(t, err)
	assert.Empty(t, fetchedUsers) //fails because delete is not fixed yet
}
