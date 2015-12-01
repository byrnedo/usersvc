package models
import (
	"testing"
	"github.com/byrnedo/usersvc/msgspec"
	"reflect"
	"github.com/byrnedo/apibase/db/mongo"
	. "github.com/byrnedo/apibase/logger"
)


func TestMain(m *testing.M){

	InitLog(func(o *LogOptions){ o.Level = InfoLevel})
	mongo.Init("mongodb://localhost:27017/test_usersvc", Info)

	m.Run()

	c := mongo.Conn()
	defer c.Close()

	c.DB("test_usersvc").DropDatabase()
}

func TestCreateUser(t *testing.T) {
	m := DefaultUserModel{}

	user, err := m.Create(&msgspec.NewUser{
		FirstName: "Test",
		LastName: "User",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}

	col, ses := m.getSession()
	defer ses.Close()

	foundUser := msgspec.UserEntity{}
	err = col.FindId(user.ID).One(&foundUser)
	if err != nil {
		t.Error("Failed to find:" + err.Error())
	}

	if reflect.DeepEqual(user, &foundUser) == false {
		t.Errorf("Did not match\nexpected:%+v\n   found:%+v\n", user, foundUser)
	}

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
	if err != nil {
		t.Error("Failed to find:" + err.Error())
	}

	if reflect.DeepEqual(user, foundUser) == false {
		t.Errorf("Did not match\nexpected:%+v\n   found:%+v\n", user, foundUser)
	}
}

func TestDeleteUser(t *testing.T) {
	m := DefaultUserModel{}

	user, err := m.Create(&msgspec.NewUser{
		FirstName: "Test",
		LastName: "User",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}

	err = m.Delete(user.ID)
	if err != nil {
		t.Error("Failed to delete:" + err.Error())
	}

}
