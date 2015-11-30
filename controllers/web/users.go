package web

import (
	routes "github.com/byrnedo/apibase/routes"
	"net/http"
)

type UsersController struct {
}

func (pC *UsersController) GetRoutes() []*routes.WebRoute{
	return []*routes.WebRoute{
		routes.NewWebRoute("GetUsers", "/api/v1/users", routes.GET, pC.List),
	}
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
