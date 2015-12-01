package models
import (
	"testing"
	"github.com/byrnedo/usersvc/msgspec"
	"reflect"
	"github.com/byrnedo/apibase/db/mongo"
	. "github.com/byrnedo/apibase/logger"
)

func init(){
	InitLog(&LogOptions{})
	mongo.Init("mongodb://localhost:27017/test_usersvc", Info)
}

func TestFindUser(t *testing.T) {

	m := DefaultUserModel{}

	user, err := m.Create(&msgspec.NewUser{
		FirstName: "Test",
		LastName: "User",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}

	foundUser, err := m.Find(user.ID)

	if reflect.DeepEqual(user, foundUser) == false {
		t.Error("Did not match")
	}
}
