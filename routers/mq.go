package routers

import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/natsio/defaultnats"
	"github.com/byrnedo/usersvc/controllers/mq"
)

func init() {
	controllers.SubscribeNatsRoutes(defaultnats.Conn, "user_svc_worker", mq.NewHealthcheckController(defaultnats.Conn))
	controllers.SubscribeNatsRoutes(defaultnats.Conn, "user_svc_worker", mq.NewUsersController(defaultnats.Conn))
}
