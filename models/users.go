package models

import (
	"github.com/byrnedo/apibase/helpers/strings"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/usersvc/msgspec"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/apibase/config"
	encBson "github.com/maxwellhealth/encryptedbson"
	"github.com/byrnedo/apibase/db/mongo/defaultmongo"
)

const (
	collection       = "users"
	defaultUserEmail = "admin@apibase.com"
)

type UserModel interface {
	Find(bson.ObjectId) (*msgspec.UserEntity, error)
	FindMany(query map[string]interface{}, sortBy []string, offset int, limit int) ([]*msgspec.UserEntity, error)
	Create(*msgspec.NewUserDTO) (*msgspec.UserEntity, error)
	Replace(*msgspec.UpdateUserDTO) (*msgspec.UserEntity, error)
	Authenticate(email string, password string) *AuthenticationError
	Delete(bson.ObjectId) error
}

type DefaultUserModel struct {
	Session *mgo.Session
}

func init(){

	encryptionKey, err := config.Conf.GetString("encryption-key")
	if err != nil {
		panic("Failed to get encryption-key:" + err.Error())
	}
	copy(encBson.EncryptionKey[:], encryptionKey)

	userModel := NewDefaultUserModel()
	userModel.Ensures()
}

func NewDefaultUserModel() *DefaultUserModel {
	return &DefaultUserModel{defaultmongo.Conn()}
}

func (u *DefaultUserModel) Ensures() {
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: false, // See notes.
		Sparse:     true,
	}
	if err := u.col().EnsureIndex(index); err != nil {
		panic("Failed to create index:" + err.Error())
	}

	if _, err := u.FindByEmail(defaultUserEmail); err != nil {
		Error.Println(err)
		var (
			randomPass = strings.RandString(12)
		)
		Info.Println("Creating default user : "+defaultUserEmail, ",  password : "+randomPass)
		if _, err = u.Create(&msgspec.NewUserDTO{
			Alias:     "defaultuser",
			FirstName: "Admin",
			LastName:  "User",
			Email:     defaultUserEmail,
			Password:  randomPass,
		}); err != nil {
			panic("Failed to create default user:" + err.Error())
		}
	}
}

func (uM *DefaultUserModel) col() *mgo.Collection {
	return uM.Session.DB("").C(collection)
}

func (uM *DefaultUserModel) Find(id bson.ObjectId) (u *msgspec.UserEntity, err error) {
	u = &msgspec.UserEntity{}
	q := uM.col().FindId(id).One(u)
	return u, q
}

func (uM *DefaultUserModel) FindByEmail(email string) (u *msgspec.UserEntity, err error) {
	u = &msgspec.UserEntity{}
	q := uM.col().Find(bson.M{"email": email}).One(u)
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
		Update:    bson.M{"$set": u},
		ReturnNew: true,
	}
	_, err = uM.col().Find(bson.M{"_id": id}).Apply(change, u)
	return
}

func (uM *DefaultUserModel) Delete(id bson.ObjectId) error {
	return uM.col().RemoveId(id)
}

func (uM *DefaultUserModel) FindMany(query map[string]interface{}, sortBy []string, offset int, limit int) ([]*msgspec.UserEntity, error) {
	var (
		err    error
		result = make([]*msgspec.UserEntity, 0)
	)

	Error.Printf("%+v", query)
	err = uM.col().Find(query).Skip(offset).Limit(limit).Sort(sortBy...).All(&result)
	return result, err
}

type AuthenticationErrorStatus string

const (
	USER_NOT_FOUND  AuthenticationErrorStatus = "User not found."
	PASSWORD_FAILED AuthenticationErrorStatus = "Password doesn't match."
)

type AuthenticationError struct {
	Reason AuthenticationErrorStatus
}

func (a *AuthenticationError) Error() string {
	return string(a.Reason)
}

func (uM *DefaultUserModel) Authenticate(email string, password string) (retErr *AuthenticationError) {
	var (
		user = &msgspec.UserEntity{}
		err  error
	)
	if user, err = uM.FindByEmail(email); err != nil {
		Error.Println(err)
		retErr = &AuthenticationError{USER_NOT_FOUND}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		Error.Println(err)
		retErr = &AuthenticationError{PASSWORD_FAILED}
	}
	return retErr
}
