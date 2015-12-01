package models
import (
	"testing"
	"github.com/byrnedo/usersvc/msgspec"
	"reflect"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/db/mongo"
)

func TestMain(m *testing.M){

	InitLog(func(o *LogOptions){ o.Level = InfoLevel})
	mongo.Init("mongodb://localhost:27017/test_mongo_model", Trace)

	m.Run()

	c := mongo.Conn()
	defer c.Close()

	c.DB("test_mongo_model").DropDatabase()
}

func TestInsertUpdateDelete(t *testing.T) {
	m := NewDefaultUserModel()
	defer m.Session.Close()

	user, err := m.Create(&msgspec.NewUser{
		FirstName: "test",
		LastName: "user",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}

	foundUser, err := m.Find(user.ID)

	if reflect.DeepEqual(user, foundUser) == false {
		t.Error("Did not match\nexpected:%+v\n   found:%+v\n",user, foundUser)
	}

	updUser, err := m.Replace(&msgspec.UpdateUser{
		ID: user.ID.Hex(),
		FirstName: "test",
		LastName: "user",
		Alias:"testy",
	})

	if err != nil {
		t.Error("Failed to update:" + err.Error())
	}

	if updUser.CreationTime != user.CreationTime {
		t.Errorf("Creation time has changed %s -> %s", user.CreationTime, updUser.CreationTime)
	}

	if updUser.UpdateTime.After(user.UpdateTime) == false {
		t.Errorf("Update time is not after insert time: %s -> %s", user.UpdateTime, updUser.UpdateTime)
	}

	err = m.Delete(user.ID)
	if err !=  nil {
		t.Error("Failed to delete user:" +err.Error())
	}
}

func TestFind(t *testing.T){

}
