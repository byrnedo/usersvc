package routers

import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/usersvc/controllers/mq"
)

func InitMq(natsCon *natsio.Nats) {
	controllers.SubscribeNatsRoutes(natsCon, "user_svc_worker", mq.NewHealthcheckController(natsCon))
	controllers.SubscribeNatsRoutes(natsCon, "user_svc_worker", mq.NewUsersController(natsCon))
}
