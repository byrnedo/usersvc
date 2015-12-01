package models
import (

	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec"
	"github.com/byrnedo/apibase/db/mongo"
	"gopkg.in/mgo.v2"
)

const (
	collection = "users"
)

type UserModel interface {
	Find(bson.ObjectId) (*msgspec.UserEntity, error)
	FindMany(query map[string]string, sortBy[]string, offset int, limit int) ([]*msgspec.UserEntity, error)
	Create(*msgspec.NewUser) (*msgspec.UserEntity, error)
	Replace(*msgspec.UpdateUser) (*msgspec.UserEntity, error)
	Delete(bson.ObjectId) error
}

type DefaultUserModel struct {
	Session *mgo.Session
}

func NewDefaultUserModel() *DefaultUserModel {
	return &DefaultUserModel{mongo.Conn()}
}

func (uM *DefaultUserModel) col() *mgo.Collection{
	return uM.Session.DB("").C(collection)
}

func (uM *DefaultUserModel) Find(id bson.ObjectId) (u *msgspec.UserEntity, err error) {
	u = &msgspec.UserEntity{}
	q := uM.col().FindId(id).One(u)
	return u, q
}

func (uM *DefaultUserModel) Create(nUser *msgspec.NewUser) (u *msgspec.UserEntity, err error) {
	u = nUser.MapToEntity()
	return u, uM.col().Insert(u)
}


func (uM *DefaultUserModel) Replace(updUser *msgspec.UpdateUser) (u *msgspec.UserEntity, err error) {
	u = updUser.MapToEntity()
	var id = u.ID
	u.ID = ""

	change := mgo.Change{
		Update: bson.M{"$set": u},
		ReturnNew: true,
	}
	_, err = uM.col().Find(bson.M{"_id": id}).Apply(change,u)
	return
}

func (uM *DefaultUserModel) Delete(id bson.ObjectId) error {
	return uM.col().RemoveId(id)
}

func (uM *DefaultUserModel) FindMany(query map[string]string, sortBy []string, offset int, limit int) ( []*msgspec.UserEntity, error ) {
	var (
		err error
		result = make([]*msgspec.UserEntity, 0)

	)
	err = mongo.GetAll(uM.col(), query, []string{}, sortBy, offset, limit, result)
	return result, err
}