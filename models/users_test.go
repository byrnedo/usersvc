package models
import (
	"testing"
	"github.com/byrnedo/usersvc/msgspec"
	"reflect"
)

func TestFindUser(t *testing.T) {
	m := DefaultUserModel{}

	user, err := m.Create(&msgspec.NewUser{

	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}

	foundUser, err := m.Find(user.ID)

	if reflect.DeepEqual(user, foundUser) == false {
		t.Error("Did not match")
	}
}
