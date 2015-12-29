package webcontrollers

import (
	"encoding/json"
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/db/mongo/defaultmongo"
	. "github.com/byrnedo/apibase/logger"
	routes "github.com/byrnedo/apibase/routes"
	svcSpec "github.com/byrnedo/svccommon/msgspec/web"
	"github.com/byrnedo/svccommon/validate"
	"github.com/byrnedo/usersvc/daos"
	"github.com/byrnedo/usersvc/models"
	"github.com/byrnedo/usersvc/msgspec/webmsgspec"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UsersController struct {
	*controllers.JsonController
	userModel daos.UserDAO
}

func NewUsersController(encryptionKey string) *UsersController {
	return &UsersController{
		JsonController: &controllers.JsonController{},
		userModel:      daos.NewDefaultUserDAO(defaultmongo.Conn(), encryptionKey), // mongo user model
	}
}

func (pC *UsersController) GetRoutes() []*routes.WebRoute {
	return []*routes.WebRoute{
		routes.NewWebRoute("CreateUser", "/v1/users", routes.POST, pC.Create),
		routes.NewWebRoute("ReplaceUser", "/v1/users/:userId", routes.PUT, pC.Replace),
		routes.NewWebRoute("GetUser", "/v1/users/:userId", routes.GET, pC.GetOne),
		routes.NewWebRoute("GetUsers", "/v1/users", routes.GET, pC.List),
		routes.NewWebRoute("DeleteUser", "/v1/users/:userId", routes.DELETE, pC.Delete),
	}
}

func (pC *UsersController) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var u webmsgspec.NewUserResource

	if err := decoder.Decode(&u); err != nil {
		Error.Println(err)
		panic("Failed to decode json:" + err.Error())
	}

	if valErrs := validate.ValidateStruct(u); len(valErrs) != 0 {
		errResponse := svcSpec.NewValidationErrorResonse(valErrs)
		pC.ServeWithStatus(w, errResponse, 400)
		return
	}

	inserted, err := pC.userModel.Create(u.Data)
	if err != nil {
		Error.Println("Error creating user:" + err.Error())
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(500), 500)
		return
	}
	pC.ServeWithStatus(w, inserted, 201)
}

func (pC *UsersController) Replace(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("userId")
	if !bson.IsObjectIdHex(id) {
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var u webmsgspec.UpdatedUserResource

	if err := decoder.Decode(&u); err != nil {
		Error.Println(err)
		panic("Failed to decode json:" + err.Error())
	}

	u.Data.ID = id

	if valErrs := validate.ValidateStruct(u); len(valErrs) != 0 {
		errResponse := svcSpec.NewValidationErrorResonse(valErrs)
		pC.ServeWithStatus(w, errResponse, 400)
		return
	}

	inserted, err := pC.userModel.Replace(u.Data)
	if err != nil {
		Error.Println("Error updating user:" + err.Error())
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(500), 500)
		return
	}
	pC.ServeWithStatus(w, inserted, 200)
}

func (pC *UsersController) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		id    string
		err   error
		objId bson.ObjectId
		user  *models.UserModel
	)

	id = ps.ByName("userId")
	if !bson.IsObjectIdHex(id) {
		Error.Println("Id is not object id")
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	objId = bson.ObjectIdHex(id)

	if user, err = pC.userModel.Find(objId); err != nil {
		Error.Println("Failed to find user:" + err.Error())
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	pC.Serve(w, &webmsgspec.UserResource{user})
}

func (pC *UsersController) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := pC.QueryInterfaceMap(r, "query", &models.UserModel{})
	order, _ := r.URL.Query()["order"]
	offset, _ := pC.QueryInt(r, "offset")
	limit, _ := pC.QueryInt(r, "limit")

	users, err := pC.userModel.FindMany(query, order, offset, limit)
	if err != nil {
		Error.Println("Failed to find users:", err)
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}
	pC.Serve(w, &webmsgspec.UsersResource{users})
}

func (pC *UsersController) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("userId")

	if !bson.IsObjectIdHex(id) {
		Error.Println("Not an object id:", id)
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	if err := pC.userModel.Delete(bson.ObjectIdHex(id)); err != nil {
		Error.Println("Error deleting:", err)
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	pC.Serve(w, &webmsgspec.UsersResource{nil})
}
