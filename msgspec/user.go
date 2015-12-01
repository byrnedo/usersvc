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
	var now = bson.Now()
	return &UserEntity{
		ID: bson.NewObjectId(),
		Alias: nU.Alias,
		FirstName: nU.FirstName,
		LastName: nU.LastName,
		Email: nU.Email,
		Password: nU.Password,
		Role: nU.Role,
		CreationTime: now,
		UpdateTime: now,
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
	UpdateTime time.Time
}

func(uU *UpdateUser) MapToEntity() (u *UserEntity) {
	var now = bson.Now()
	return &UserEntity{
		Alias: uU.Alias,
		FirstName: uU.FirstName,
		LastName: uU.LastName,
		Email: uU.Email,
		Password: uU.Password,
		Role: uU.Role,
		UpdateTime: now,
	}
}

func (u *UpdateUser) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}


type UserEntity struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Alias string
	FirstName string
	LastName string
	Email string
	Password string `bson:"omitempty"`
	Role string
	CreationTime time.Time
	UpdateTime time.Time
}


