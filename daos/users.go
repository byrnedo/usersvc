package daos
import (
	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec/webmsgspec"
	"gopkg.in/mgo.v2"
	"github.com/byrnedo/usersvc/models"
	"github.com/byrnedo/apibase/helpers/stringhelp"
	encBson "github.com/maxwellhealth/encryptedbson"
	. "github.com/byrnedo/apibase/logger"
	"golang.org/x/crypto/bcrypt"
)

const (
	collection       = "users"
	defaultUserEmail = "admin@apibase.com"
)


type UserDAO interface {
	Find(bson.ObjectId) (*models.UserModel, error)
	FindMany(query map[string]interface{}, sortBy []string, offset int, limit int) ([]*models.UserModel, error)
	Create(*webmsgspec.NewUserDTO) (*models.UserModel, error)
	Replace(*webmsgspec.UpdateUserDTO) (*models.UserModel, error)
	Authenticate(email string, password string) *AuthenticationError
	Delete(bson.ObjectId) error
}

type DefaultUserDAO struct {
	Session *mgo.Session
}

func init(){
}

func NewDefaultUserDAO(session *mgo.Session, encryptionKey string) *DefaultUserDAO {
	copy(encBson.EncryptionKey[:], encryptionKey)

	dao := &DefaultUserDAO{session}
	dao.Ensures()
	return dao
}

func (u *DefaultUserDAO) Ensures() {
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
			randomPass = stringhelp.RandString(12)
		)
		Info.Println("Creating default user : "+defaultUserEmail, ",  password : "+randomPass)
		if _, err = u.Create(&webmsgspec.NewUserDTO{
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

func (uM *DefaultUserDAO) col() *mgo.Collection {
	return uM.Session.DB("").C(collection)
}

func (uM *DefaultUserDAO) Find(id bson.ObjectId) (u *models.UserModel, err error) {
	u = &models.UserModel{}
	q := uM.col().FindId(id).One(u)
	return u, q
}

func (uM *DefaultUserDAO) FindByEmail(email string) (u *models.UserModel, err error) {
	u = &models.UserModel{}
	q := uM.col().Find(bson.M{"email": email}).One(u)
	return u, q
}

func (uM *DefaultUserDAO) Create(nUser *webmsgspec.NewUserDTO) (u *models.UserModel, err error) {
	if u, err = nUser.MapToEntity(); err != nil {
		return
	}

	return u, uM.col().Insert(u)
}

func (uM *DefaultUserDAO) Replace(updUser *webmsgspec.UpdateUserDTO) (u *models.UserModel, err error) {
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

func (uM *DefaultUserDAO) Delete(id bson.ObjectId) error {
	return uM.col().RemoveId(id)
}

func (uM *DefaultUserDAO) FindMany(query map[string]interface{}, sortBy []string, offset int, limit int) ([]*models.UserModel, error) {
	var (
		err    error
		result = make([]*models.UserModel, 0)
	)

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

func (uM *DefaultUserDAO) Authenticate(email string, password string) (retErr *AuthenticationError) {
	var (
		user = &models.UserModel{}
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
