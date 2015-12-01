package models
import (

	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec"
	"github.com/byrnedo/apibase/db/mongo"
	"gopkg.in/mgo.v2"
)

type UserModel interface {
	Find(bson.ObjectId) (*msgspec.UserEntity, error)
	FindMany() ([]*msgspec.UserEntity, error)
	Create(*msgspec.NewUser) (*msgspec.UserEntity, error)
	Replace(*msgspec.UpdateUser) (*msgspec.UserEntity, error)
	Delete(bson.ObjectId) error
}

type DefaultUserModel struct {
}

func (uM *DefaultUserModel) getSession() (*mgo.Collection, *mgo.Session) {
	mSess := mongo.Conn()
	return mSess.DB("").C("users"), mSess
}

func (uM *DefaultUserModel) Find(id bson.ObjectId) (u *msgspec.UserEntity, err error) {
	col, ses := uM.getSession()
	defer ses.Close()

	u = &msgspec.UserEntity{}
	q := col.FindId(id).One(u)
	return u, q
}

func (uM *DefaultUserModel) Create(nUser *msgspec.NewUser) (u *msgspec.UserEntity, err error) {
	col, ses := uM.getSession()
	defer ses.Close()

	u = nUser.MapToEntity()

	return u, col.Insert(u)
}
