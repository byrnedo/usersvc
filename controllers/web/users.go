package web

import (
	routes "github.com/byrnedo/apibase/routes"
	"net/http"
	"github.com/byrnedo/usersvc/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/byrnedo/usersvc/msgspec"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/controllers"
	svcSpec "github.com/byrnedo/svccommon/msgspec/web"
	"encoding/json"
	"github.com/byrnedo/usersvc/msgspec/web"
	"github.com/julienschmidt/httprouter"
)

type UsersController struct {
	*controllers.JsonController
	userModel models.UserModel
}

func NewUsersController() *UsersController {
	return &UsersController{
		JsonController: &controllers.JsonController{},
		userModel: models.NewDefaultUserModel(), // mongo user model
	}
}

func (pC *UsersController) GetRoutes() []*routes.WebRoute{
	return []*routes.WebRoute{
		routes.NewWebRoute("CreateUser", "/api/v1/users", routes.POST, pC.Create),
		routes.NewWebRoute("GetUser", "/api/v1/users/:userId", routes.GET, pC.GetOne),
		routes.NewWebRoute("GetUsers", "/api/v1/users", routes.GET, pC.List),
	}
}

func (pC *UsersController) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	decoder := json.NewDecoder(r.Body)
	var u web.NewUserResource

	if err := decoder.Decode(&u); err != nil {
		Error.Println(err)
		panic("Failed to decode json:"+err.Error())
	}

	if valErrs := u.Validate(); len(valErrs) != 0 {
		errResponse := svcSpec.NewErrorResponse()
		for field, fieldErr := range valErrs {
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

func (pC *UsersController) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var (
		id string
		err error
		objId bson.ObjectId
		user *msgspec.UserEntity
	)

	id = ps.ByName("userId")

	if bson.IsObjectIdHex(id) == false {
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

	pC.Serve(w, &web.UserResource{user})

}

func (pC *UsersController) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params){

	//var query = pC.QueryMap(r, "query")
	//var fields = pC.QuerySlice("fields")
	var order = pC.QuerySlice(r, "order")


	offset, _ := pC.QueryInt(r, "offset")
	limit, _ := pC.QueryInt(r, "limit")

	users, err := pC.userModel.FindMany(nil,order,offset,limit)
	if err != nil {
		Error.Println("Failed to find users:" + err.Error())
		pC.ServeWithStatus(w, svcSpec.NewErrorResponse().AddCodeError(404), 404)
		return
	}

	pC.Serve(w, &web.ManyUserResource{users})

}
