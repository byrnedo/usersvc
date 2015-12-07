package web

import (
	routes "github.com/byrnedo/apibase/routes"
	"net/http"
	"github.com/fsouza/go-dockerclient/external/github.com/gorilla/mux"
	"github.com/byrnedo/usersvc/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec"
)

type UsersController struct {
	userModel models.UserModel
}

func NewUsersController() *UsersController {
	return &UsersController{
		userModel: &models.DefaultUserModel{},
	}
}

func (pC *UsersController) GetRoutes() []*routes.WebRoute{
	return []*routes.WebRoute{
		routes.NewWebRoute("GetUsers", "/api/v1/users", routes.GET, pC.List),
		routes.NewWebRoute("GetUsers", "/api/v1/users/{id}", routes.GET, pC.GetOne),
	}
}

func (pC *UsersController) GetOne(w http.ResponseWriter, r *http.Request){
	var (
		id string
		err error
		found bool
		objId bson.ObjectId
		user *msgspec.UserEntity
	)
	if id, found = mux.Vars(r)["id"]; found == false {
		//send error
		return
	}

	if bson.IsObjectIdHex(id) == false {
		//send error
		return
	}

	objId = bson.ObjectIdHex(id)


	if user, err = pC.userModel.Find(objId); err != nil {
		//send 404
		return
	}

	_ = user
	//send json


}

func (pC *UsersController) List(w http.ResponseWriter, r *http.Request){

//	var fields = r.URL.Query("fields")
//
//	var order = r.URL.Query("order")
//
//	var offset = r.URL.Query("offset")
//	var limit = r.URL.Query("limit")
//
//	r.URL.Query().Get()

}
