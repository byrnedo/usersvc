package models
import (
	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec"
	"github.com/byrnedo/apibase/db/mongo"
	"gopkg.in/mgo.v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	collection = "users"
)

type UserModel interface {
	Find(bson.ObjectId) (*msgspec.UserEntity, error)
	FindMany(query map[string]interface{}, sortBy[]string, offset int, limit int) ([]*msgspec.UserEntity, error)
	Create(*msgspec.NewUserDTO) (*msgspec.UserEntity, error)
	Replace(*msgspec.UpdateUserDTO) (*msgspec.UserEntity, error)
	Authenticate(email string, password string) bool
	Delete(bson.ObjectId) error
}

type DefaultUserModel struct {
	Session *mgo.Session
}

func init() {

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

func (uM *DefaultUserModel) Create(nUser *msgspec.NewUserDTO) (u *msgspec.UserEntity, err error) {
	if u, err = nUser.MapToEntity(); err != nil {
		return
	}

	return u, uM.col().Insert(u)
}


func (uM *DefaultUserModel) Replace(updUser *msgspec.UpdateUserDTO) (u *msgspec.UserEntity, err error) {
	if u, err = updUser.MapToEntity(); err != nil {
		return
	}
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

func (uM *DefaultUserModel) FindMany(query map[string]interface{}, sortBy []string, offset int, limit int) ( []*msgspec.UserEntity, error ) {
	var (
		err error
		result = make([]*msgspec.UserEntity, 0)

	)
	mongo.ConvertObjectIds(query)
	err = uM.col().Find(query).Skip(offset).Limit(limit).Sort(sortBy...).All(&result)
	return result, err
}

func (uM *DefaultUserModel) Authenticate(email string, password string) bool {
	var user = &msgspec.UserEntity{}
	if err := uM.col().Find(bson.M{"email": email}).One(user); err != nil {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}