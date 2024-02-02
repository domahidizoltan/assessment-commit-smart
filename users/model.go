package users

import (
	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            *string
	Id               int
	Password         *string
	Username         *string
}
