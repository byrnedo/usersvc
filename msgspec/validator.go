package msgspec

import (
	validator "gopkg.in/bluesuncorp/validator.v8"
	"reflect"
)

var V *validator.Validate

func init() {
	V = validator.New(&validator.Config{TagName: "validate"})
}

func ValidateStruct(any interface{}) (valErrors map[string]*validator.FieldError) {

	var anyKind reflect.Kind = reflect.TypeOf(any).Kind()

	if anyKind != reflect.Struct {
		panic("can only pass structs to this function, instead got " + anyKind.String() )
	}

	if errs := V.Struct(any); errs != nil {
		valErrors = errs.(validator.ValidationErrors)
	}
	return
}
