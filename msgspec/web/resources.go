package web

import (
	"github.com/byrnedo/usersvc/msgspec"
	"gopkg.in/bluesuncorp/validator.v8"
	"github.com/byrnedo/svccommon/validate"
)

type NewUserResource struct {
	Data *msgspec.NewUserDTO `json:"data" validate:"required"`
}

type UpdatedUserResource struct {
	Data *msgspec.UpdateUserDTO `json:"data" validate:"required"`
}

type UserResource struct {
	Data *msgspec.UserEntity `json:"data"`
}

type UsersResource struct {
	Data []*msgspec.UserEntity `json:"data"`
}

func (nU *NewUserResource) Validate() validator.ValidationErrors {
	return validate.ValidateStruct(nU)
}
