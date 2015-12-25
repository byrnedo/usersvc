package models

import (
	"github.com/byrnedo/svccommon/validate"
	encBson "github.com/maxwellhealth/encryptedbson"
	"gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type RoleType string

const (
	AdminRole  RoleType = "admin"
	NormalRole RoleType = "normal"
)

type UserModel struct {
	ID           bson.ObjectId           `bson:"_id,omitempty" json:"id"`
	Alias        string                  `json:"alias"`
	FirstName    encBson.EncryptedString `json:"firstname"`
	LastName     encBson.EncryptedString `json:"lastname"`
	Email        string                  `json:"email"`
	Password     string                  `bson:"password,omitempty" json:"-"`
	Role         RoleType                `json:"role"`
	CreationTime time.Time               `bson:"creationtime,omitempty" json:"creationtime"`
	UpdateTime   time.Time               `json:"updatetime"`
}

func (u *UserModel) Validate() validator.ValidationErrors {
	return validate.ValidateStruct(u)
}
