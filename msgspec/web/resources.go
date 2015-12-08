package web
import (
	"github.com/byrnedo/usersvc/msgspec"
	"gopkg.in/bluesuncorp/validator.v8"
	"github.com/byrnedo/apibase/validate"
)

type NewUserResource struct {
	Data *msgspec.NewUserDTO `json:"data" validate:"required"`
}

func (nU *NewUserResource) Validate() validator.ValidationErrors {
	return validate.ValidateStruct(nU)
}
