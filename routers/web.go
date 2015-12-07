package routers

import (
	"github.com/gorilla/mux"
	"github.com/byrnedo/usersvc/controllers/web"
	"github.com/byrnedo/apibase/controllers"
	"net/http"
	"github.com/justinas/alice"
	//"github.com/ulule/limiter"
	"github.com/byrnedo/apibase/middleware"
)

func InitWeb() {
	var rtr = mux.NewRouter().StrictSlash(true)
	controllers.RegisterMuxRoutes(rtr, web.NewUsersController())

	//alice is a tiny package to chain middlewares.
	mChain := alice.New(middleware.LogTime).Then(rtr)

	http.Handle("/", mChain)
}

