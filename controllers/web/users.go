package web

import (
	routes "github.com/byrnedo/apibase/routes"
	"net/http"
	"github.com/fsouza/go-dockerclient/external/github.com/gorilla/mux"
	"github.com/byrnedo/usersvc/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/controllers"
	svcSpec "github.com/byrnedo/svccommon/msgspec/web"
	"encoding/json"
	"github.com/byrnedo/usersvc/msgspec/web"
)

type UsersController struct {
	*controllers.JsonController
	userModel models.UserModel
}

func NewUsersController() *UsersController {
	return &UsersController{
		JsonController: &controllers.JsonController{},
		userModel: &models.DefaultUserModel{}, // mongo user model
	}
}

func (pC *UsersController) GetRoutes() []*routes.WebRoute{
	return []*routes.WebRoute{
		routes.NewWebRoute("CreateUser", "/api/v1/users", routes.POST, pC.Create),
		routes.NewWebRoute("GetUsers", "/api/v1/users", routes.GET, pC.List),
		routes.NewWebRoute("GetUser", "/api/v1/users/{id}", routes.GET, pC.GetOne),
	}
}

func (pC *UsersController) Create(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var u web.NewUserResource
	err := decoder.Decode(&u)
	if err != nil {
		Error.Println(err)
		panic("Failed to decode json:"+err.Error())
	}

	if valErrs := u.Validate(); len(valErrs) != 0 {
		errResponse := svcSpec.NewErrorResponse()

		for field, fieldErr := range valErrs {
			Info.Printf("%s -> %+v\n", field, fieldErr)
			errResponse.AddError(400, &svcSpec.Source{Parameter:field}, fieldErr.Tag, fieldErr.Tag)
		}
		pC.ServeWithStatus(w, errResponse, 400)
		return
	}

	inserted, err := pC.userModel.Create(u.Data)
	if err != nil {
		Error.Println("Error creating user:"+err.Error())
		pC.ServeWithStatus(w,svcSpec.NewErrorResponse().AddCodeError(500), 500)
		return
	}
	pC.ServeWithStatus(w, inserted, 201)
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
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	if bson.IsObjectIdHex(id) == false {
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	objId = bson.ObjectIdHex(id)


	if user, err = pC.userModel.Find(objId); err != nil {
		Error.Println("Failed to find user:" + err.Error())
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	pC.Serve(w, user)

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
