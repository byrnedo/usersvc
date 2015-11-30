package msgspec

import (
	"time"
	validator "gopkg.in/bluesuncorp/validator.v8"
)

type NewUser struct {
	User
	Password string `validate:"required"`
}

type UpdateUser struct {
	ID string `validate:"required"`
	User
}


type User struct {
	Alias string `validate:"required"`
	FirstName string
	LastName string
	Email string `validate:"required"`
	Password string
	Role string `validate:"required"`
	CreationTime time.Time
	UpdateTime time.Time
}

func (u *User) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}

