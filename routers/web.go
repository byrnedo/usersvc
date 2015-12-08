package routers

import (
	"github.com/gorilla/mux"
	"github.com/byrnedo/usersvc/controllers/web"
	"github.com/byrnedo/apibase/controllers"
	"net/http"
	"github.com/justinas/alice"
	"github.com/ulule/limiter"
	"github.com/byrnedo/apibase/middleware"
	"time"
)

func InitWeb() {

	rate, err := limiter.NewRateFromFormatted("5-S")
	if err != nil {
		panic(err)
	}

	//TODO - exchange for redis store
	store := limiter.NewMemoryStoreWithOptions(limiter.StoreOptions{
		Prefix:"byrnedosvc",
		CleanUpInterval: 30*time.Second,
	})
	limiterMw := limiter.NewHTTPMiddleware(limiter.NewLimiter(store, rate))


	var rtr = mux.NewRouter().StrictSlash(true)
	controllers.RegisterMuxRoutes(rtr, web.NewUsersController())

	//alice is a tiny package to chain middlewares.
	handlerChain := alice.New(
		limiterMw.Handler,
		middleware.LogTime,
		middleware.RecoverHandler,
		middleware.AcceptJsonHandler,
	).Then(rtr)

	http.Handle("/", handlerChain)
}

