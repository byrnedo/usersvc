package webmsgspec

import (
	"errors"
	"github.com/byrnedo/usersvc/models"
	encBson "github.com/maxwellhealth/encryptedbson"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const (
	bcryptCost = 10
)

type NewUserResource struct {
	Data *NewUserDTO `json:"data" validate:"required"`
}

type UpdatedUserResource struct {
	Data *UpdateUserDTO `json:"data" validate:"required"`
}

type UserResource struct {
	Data *models.UserModel `json:"data"`
}

type UsersResource struct {
	Data []*models.UserModel `json:"data"`
}

type NewUserDTO struct {
	Alias     string          `json:"alias" validate:"required,alphanum"`
	FirstName string          `json:"first_name" validate:"omitempty,alpha"`
	LastName  string          `json:"last_name" validate:"omitempty,alpha"`
	Email     string          `json:"email" validate:"required,email"`
	Password  string          `json:"password" validate:"omitempty,min=8"`
	Role      models.RoleType `json:"role" validate:"required,eq=admin|eq=normal"`
}

func (nU *NewUserDTO) MapToEntity() (*models.UserModel, error) {
	var (
		now = bson.Now()
		err error
	)

	if nU.Password, err = encryptPassword(nU.Password); err != nil {
		return nil, err
	}

	return &models.UserModel{
		ID:           bson.NewObjectId(),
		Alias:        nU.Alias,
		FirstName:    encBson.EncryptedString(nU.FirstName),
		LastName:     encBson.EncryptedString(nU.LastName),
		Email:        nU.Email,
		Password:     nU.Password,
		Role:         nU.Role,
		CreationTime: now,
		UpdateTime:   now,
	}, nil
}

type UpdateUserDTO struct {
	ID        string          `json:"id"`
	Alias     string          `json:"alias" validate:"required,alphanum"`
	FirstName string          `json:"first_name" validate:"omitempty,alpha"`
	LastName  string          `json:"last_name" validate:"omitempty,alpha"`
	Email     string          `json:"email" validate:"required,email"`
	Password  string          `json:"password" validate:"omitempty,min=8"`
	Role      models.RoleType `json:"role" validate:"required,eq=admin|eq=normal"`
}

func (uU *UpdateUserDTO) MapToEntity() (*models.UserModel, error) {

	var err error
	if len(uU.Password) > 0 {
		if uU.Password, err = encryptPassword(uU.Password); err != nil {
			return nil, err
		}
	}
	return &models.UserModel{
		ID:         bson.ObjectIdHex(uU.ID),
		Alias:      uU.Alias,
		FirstName:  encBson.EncryptedString(uU.FirstName),
		LastName:   encBson.EncryptedString(uU.LastName),
		Email:      uU.Email,
		Password:   uU.Password,
		Role:       uU.Role,
		UpdateTime: bson.Now(),
	}, nil
}

func encryptPassword(pass string) (string, error) {

	password, err := bcrypt.GenerateFromPassword([]byte(pass), bcryptCost)
	if err != nil {
		return "", errors.New("Failed to encrypt:" + err.Error())
	}
	return string(password), nil

}
