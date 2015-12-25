package daos

import (
	"github.com/byrnedo/apibase/db/mongo"
	"github.com/byrnedo/apibase/dockertest"
	"github.com/byrnedo/apibase/helpers/stringhelp"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/usersvc/models"
	"github.com/byrnedo/usersvc/msgspec/webmsgspec"
	gDoc "github.com/fsouza/go-dockerclient"
	"gopkg.in/mgo.v2"
	"reflect"
	"testing"
	"time"
)

var (
	session *mgo.Session
)

func Conn() *mgo.Session {
	return session.Copy()
}

const (
	MongoImage = "mongo:latest"
	MongoPort  = "28017"
)

func setupContainer() {

	if id, err := dockertest.Running(MongoImage); err != nil || len(id) < 1 {
		if _, err := dockertest.Start(MongoImage, map[gDoc.Port][]gDoc.PortBinding{
			"27017/tcp": []gDoc.PortBinding{gDoc.PortBinding{
				HostIP:   "127.0.0.1",
				HostPort: MongoPort,
			}},
		}); err != nil {
			panic("Error starting postgres:" + err.Error())
		}

	}
}

func TestMain(m *testing.M) {

	setupContainer()

	InitLog(func(o *LogOptions) { o.Level = InfoLevel })
	session = mongo.Init("mongodb://localhost:"+MongoPort+"/test_users", Trace)

	c := Conn()
	c.DB("test_users").DropDatabase()
	defer c.Close()
	m.Run()

}

func createUser(m UserDAO, t *testing.T) *models.UserModel {

	user, err := m.Create(&webmsgspec.NewUserDTO{
		FirstName: "test",
		LastName:  "user",
		Alias:     "test",
		Password:  "test",
		Email:     stringhelp.RandString(5) + "@apibase.com",
	})

	if err != nil {
		t.Error("Failed to insert:" + err.Error())
	}
	time.Sleep(10 * time.Millisecond)
	return user
}

func TestModel(t *testing.T) {
	m := &DefaultUserDAO{session.Copy()}
	defer m.Session.Close()

	user := createUser(m, t)

	foundUser, err := m.Find(user.ID)

	t.Logf("%+v", foundUser)

	if reflect.DeepEqual(user, foundUser) == false {
		t.Error("Did not match\nexpected:%+v\n   found:%+v\n", user, foundUser)
	}

	updUser, err := m.Replace(&webmsgspec.UpdateUserDTO{
		ID:        user.ID.Hex(),
		FirstName: "test",
		LastName:  "user",
		Email:     "email",
		Password:  "password",
		Alias:     "testy",
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

	if valid != nil {
		t.Error("Authenticate returned false")
	}

	valid = m.Authenticate("email", "password2")
	if valid == nil {
		t.Error("Bad Authenticate returned true")
	}

	valid = m.Authenticate("gmail", "password2")
	if valid == nil {
		t.Error("Bad Authenticate returned true")
	}

	err = m.Delete(user.ID)
	if err != nil {
		t.Error("Failed to delete user:" + err.Error())
	}

	var users = []*models.UserModel{
		createUser(m, t),
		createUser(m, t),
		createUser(m, t),
	}

	res, err := m.FindMany(nil, []string{"-creationtime"}, 0, 10)
	if err != nil {
		t.Error("Failed to find many:" + err.Error())
	}
	if len(res) != 3 {
		t.Errorf("Got %d results, expected 3\n", len(res))
	}

	res, err = m.FindMany(nil, []string{"-creationtime"}, 0, 1)
	if err != nil {
		t.Error("Failed to find many:" + err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Got %d results, expected 1\n", len(res))
	}

	if res[0].ID != users[2].ID {
		t.Errorf("Not the user, \nexpected: %+v\n   found: %+v", users[2], res[0])
	}

}
