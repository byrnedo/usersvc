package webcontrollers

import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() (*httptest.ResponseRecorder, *mux.Router) {

	rtr := mux.NewRouter().StrictSlash(true)
	controllers.RegisterRoutes(rtr, &UsersController{})
	rec := httptest.NewRecorder()
	return rec, rtr
}

func Ensure200(status int, t *testing.T) {
	if status != http.StatusOK {
		t.Errorf("Home page didn't return %v, got %d", http.StatusOK, status)
	}
}

func TestGetAllUsers(t *testing.T) {
	rec, rtr := setup()

	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	if err != nil {
		t.Error(err.Error())
	}
	rtr.ServeHTTP(rec, req)

	Ensure200(rec.Code, t)
}
