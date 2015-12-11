package msgspec

import (
	"errors"
	encBson "github.com/maxwellhealth/encryptedbson"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/byrnedo/svccommon/validate"
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
	Alias        string    `json:"alias"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	CreationTime time.Time `json:"creationtime"`
	UpdateTime   time.Time `json:"updatetime"`
}

func (nU *NewUserDTO) MapToEntity() (*UserEntity, error) {
	var (
		now = bson.Now()
		err error
	)

	if nU.Password, err = encryptPassword(nU.Password); err != nil {
		return nil, err
	}

	return &UserEntity{
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

func (u *NewUserDTO) Validate() validator.ValidationErrors {
	return validate.ValidateStruct(u)
}

type UpdateUserDTO struct {
	ID           string    `json:"id"`
	Alias        string    `json:"alias"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	CreationTime time.Time `json:"creationtime"`
	UpdateTime   time.Time `json:"updatetime"`
}

func (uU *UpdateUserDTO) MapToEntity() (*UserEntity, error) {

	var err error
	if len(uU.Password) > 0 {
		if uU.Password, err = encryptPassword(uU.Password); err != nil {
			return nil, err
		}
	}
	return &UserEntity{
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

func (u *UpdateUserDTO) Validate() validator.ValidationErrors {
	return validate.ValidateStruct(u)
}

type UserEntity struct {
	ID           bson.ObjectId           `bson:"_id,omitempty" json:"id"`
	Alias        string                  `json:"alias"`
	FirstName    encBson.EncryptedString `json:"firstname"`
	LastName     encBson.EncryptedString `json:"lastname"`
	Email        string                  `json:"email"`
	Password     string                  `bson:"password,omitempty" json:"-"`
	Role         string                  `json:"role"`
	CreationTime time.Time               `bson:"creationtime,omitempty" json:"creationtime"`
	UpdateTime   time.Time               `json:"updatetime"`
}

func (u *UserEntity) Validate() validator.ValidationErrors {
	return validate.ValidateStruct(u)
}
