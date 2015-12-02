package msgspec

import (
	"time"
	validator "gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
	encBson "github.com/maxwellhealth/encryptedbson"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

const (
	bcryptCost = 10
)

func encryptPassword(pass string) (string, error) {

	password, err := bcrypt.GenerateFromPassword([]byte(pass), bcryptCost)
	if err != nil {
		return "", errors.New("Failed to encrypt:" + err.Error())
	}
	return string(password), nil

}

type NewUserDTO struct {
	Alias string
	FirstName string
	LastName string
	Email string
	Password string
	Role string
	CreationTime time.Time
	UpdateTime time.Time
}

func(nU *NewUserDTO) MapToEntity() (*UserEntity, error) {
	var now = bson.Now()
	var err error

	if nU.Password, err = encryptPassword(nU.Password); err != nil {
		return nil,err
	}

	return &UserEntity{
		ID: bson.NewObjectId(),
		Alias: nU.Alias,
		FirstName: encBson.EncryptedString(nU.FirstName),
		LastName: encBson.EncryptedString(nU.LastName),
		Email: nU.Email,
		Password: nU.Password,
		Role: nU.Role,
		CreationTime: now,
		UpdateTime: now,
	}, nil
}

func (u *NewUserDTO) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}

type UpdateUserDTO struct {
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

func(uU *UpdateUserDTO) MapToEntity() (*UserEntity, error) {

	var err error
	if len(uU.Password) > 0 {
		if uU.Password, err = encryptPassword(uU.Password); err != nil {
			return nil,err
		}
	}
	return &UserEntity{
		ID: bson.ObjectIdHex(uU.ID),
		Alias: uU.Alias,
		FirstName: encBson.EncryptedString(uU.FirstName),
		LastName: encBson.EncryptedString(uU.LastName),
		Email: uU.Email,
		Password: uU.Password,
		Role: uU.Role,
		UpdateTime: bson.Now(),
	}, nil
}

func (u *UpdateUserDTO) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}


type UserEntity struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Alias string
	FirstName encBson.EncryptedString
	LastName encBson.EncryptedString
	Email string
	Password string `bson:"password,omitempty"`
	Role string
	CreationTime time.Time `bson:"creationtime,omitempty"`
	UpdateTime time.Time
}

func (u *UserEntity) Validate() map[string]*validator.FieldError {
	return V.Struct(u).(validator.ValidationErrors)
}

