package routers

import (
	"github.com/gorilla/mux"
	"github.com/byrnedo/usersvc/controllers/web"
	"github.com/byrnedo/apibase/controllers"
	"net/http"
	"github.com/justinas/alice"
	//"github.com/ulule/limiter"
)

func InitWeb() {
	var rtr = mux.NewRouter().StrictSlash(true)
	controllers.RegisterMuxRoutes(rtr, &web.UsersController{})

	//alice is a tiny package to chain middlewares.
	mChain := alice.New(Middleware).Then(rtr)

	http.Handle("/", mChain)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

