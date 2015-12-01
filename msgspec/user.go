package msgspec

import (
	"time"
	validator "gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
)

type NewUser struct {
	Alias string
	FirstName string
	LastName string
	Email string
	Password string
	Role string
	CreationTime time.Time
	UpdateTime time.Time
}

func(nU *NewUser) MapToEntity() (u *UserEntity) {
	return &UserEntity{
		Alias: nU.Alias,
		FirstName: nU.FirstName,
		LastName: nU.LastName,
		Email: nU.Email,
		Password: nU.Password,
		Role: nU.Role,
		CreationTime: nU.CreationTime,
		UpdateTime: nU.UpdateTime,
	}
}

func (u *NewUser) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}

type UpdateUser struct {
	ID string
	Alias string
	FirstName string
	LastName string
	Email string
	Password string
	Role string
	CreationTime time.Time
	UpdateTime time.Time
}

func(uU *UpdateUser) MapToEntity() (u *UserEntity) {
	return &UserEntity{
		Alias: uU.Alias,
		FirstName: uU.FirstName,
		LastName: uU.LastName,
		Email: uU.Email,
		Password: uU.Password,
		Role: uU.Role,
		CreationTime: uU.CreationTime,
		UpdateTime: uU.UpdateTime,
	}
}

func (u *UpdateUser) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}


type UserEntity struct {
	ID bson.ObjectId
	Alias string
	FirstName string
	LastName string
	Email string
	Password string
	Role string
	CreationTime time.Time
	UpdateTime time.Time
}

func (u *UserEntity) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}

