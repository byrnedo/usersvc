package models
import (
	"testing"
	"github.com/byrnedo/usersvc/msgspec"
	"reflect"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/db/mongo"
	"time"
)

func TestMain(m *testing.M){

	InitLog(func(o *LogOptions){ o.Level = InfoLevel})
	mongo.Init("mongodb://localhost:27017/test_mongo_model", Trace)

	m.Run()

	c := mongo.Conn()
	defer c.Close()

	c.DB("test_mongo_model").DropDatabase()
}

func createUser(m UserModel, t *testing.T) *msgspec.UserEntity {

	user, err := m.Create(&msgspec.NewUser{
		FirstName: "test",
		LastName: "user",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}
	time.Sleep(10*time.Millisecond)
	return user
}

func TestModel(t *testing.T) {
	m := NewDefaultUserModel()
	defer m.Session.Close()

	user := createUser(m, t)

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

	var users = []*msgspec.UserEntity {
		createUser(m, t),
		createUser(m, t),
		createUser(m, t),
	}

	res, err := m.FindMany(nil, []string{"-creationtime"}, 0,10)
	if err != nil {
		t.Error("Failed to find many:" + err.Error())
	}
	if len(res) != 3 {
		t.Errorf("Got %d results, expected 3\n", len(res))
	}

	res, err = m.FindMany(nil, []string{"-creationtime"}, 0,1)
	if err != nil {
		t.Error("Failed to find many:" + err.Error())
	}

	if len(res)!= 1 {
		t.Errorf("Got %d results, expected 1\n", len(res))
	}

	if res[0].ID != users[2].ID {
		t.Errorf("Not the user, \nexpected: %+v\n   found: %+v", users[2], res[0])
	}

}
