package routers

import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/natsio/defaultnats"
	"github.com/byrnedo/usersvc/controllers/mqcontrollers"
)

func init() {
	controllers.SubscribeNatsRoutes(defaultnats.Conn, "user_svc_worker", mqcontrollers.NewHealthcheckController(defaultnats.Conn))
	controllers.SubscribeNatsRoutes(defaultnats.Conn, "user_svc_worker", mqcontrollers.NewUsersController(defaultnats.Conn))
}
