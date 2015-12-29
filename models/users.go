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
	FirstName    encBson.EncryptedString `json:"first_name"`
	LastName     encBson.EncryptedString `json:"last_name"`
	Email        string                  `json:"email"`
	Password     string                  `bson:",omitempty" json:"-"`
	Role         RoleType                `json:"role"`
	CreationTime time.Time               `bson:",omitempty" json:"creation_time"`
	UpdateTime   time.Time               `json:"update_time"`
}
