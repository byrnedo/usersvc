package models
import (

	"testing"
	"github.com/byrnedo/usersvc/msgspec"
	"reflect"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/db/mongo"
	"time"
	"github.com/byrnedo/apibase/dockertest"
	gDoc "github.com/fsouza/go-dockerclient"
)

const(
	MongoImage = "mongo:latest"
	MongoPort = "28017"
)

func setupContainer() {

	if id, err := dockertest.Running(MongoImage); err != nil || len(id) < 1 {
		if _, err := dockertest.Start(MongoImage, map[gDoc.Port][]gDoc.PortBinding{
			"27017/tcp" : []gDoc.PortBinding{gDoc.PortBinding{
				HostIP: "127.0.0.1",
				HostPort: MongoPort,
			}},
		}); err != nil {
			panic("Error starting postgres:" + err.Error())
		}

	}
}

func TestMain(m *testing.M){

	setupContainer()

	InitLog(func(o *LogOptions){ o.Level = InfoLevel})
	mongo.Init("mongodb://localhost:"+MongoPort+"/test_mongo_model", Trace)

	c := mongo.Conn()
	c.DB("test_mongo_model").DropDatabase()
	defer c.Close()
	m.Run()


}

func createUser(m UserModel, t *testing.T) *msgspec.UserEntity {

	user, err := m.Create(&msgspec.NewUserDTO{
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

	t.Logf("%+v", foundUser)

	if reflect.DeepEqual(user, foundUser) == false {
		t.Error("Did not match\nexpected:%+v\n   found:%+v\n",user, foundUser)
	}

	updUser, err := m.Replace(&msgspec.UpdateUserDTO{
		ID: user.ID.Hex(),
		FirstName: "test",
		LastName: "user",
		Email: "email",
		Password: "password",
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


	valid := m.Authenticate("email", "password")

	if valid == false {
		t.Error("Authenticate returned false")
	}

	valid = m.Authenticate("email", "password2")
	if valid == true {
		t.Error("Bad Authenticate returned true")
	}

	valid = m.Authenticate("gmail", "password2")
	if valid == true {
		t.Error("Bad Authenticate returned true")
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
