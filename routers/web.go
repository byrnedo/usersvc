package routers

import (
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/middleware"
	"github.com/byrnedo/usersvc/controllers/webcontrollers"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func init() {

	encryptionKey, err := config.Conf.GetString("encryption-key")
	if err != nil {
		panic("Failed to get encryption-key:" + err.Error())
	}

	//	rate, err := limiter.NewRateFromFormatted("5-S")
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	//TODO - exchange for redis store
	//	store := limiter.NewMemoryStoreWithOptions(limiter.StoreOptions{
	//		Prefix:          "byrnedosvc",
	//		CleanUpInterval: 30 * time.Second,
	//	})
	//	limiterMw := limiter.NewHTTPMiddleware(limiter.NewLimiter(store, rate))

	var rtr = httprouter.New()
	controllers.RegisterRoutes(rtr, webcontrollers.NewUsersController(encryptionKey))

	//alice is a tiny package to chain middlewares.
	handlerChain := alice.New(
		//limiterMw.Handler,
		middleware.LogTime,
		middleware.RecoverHandler,
		middleware.AcceptJsonHandler,
	).Then(rtr)

	http.Handle("/", handlerChain)
}
