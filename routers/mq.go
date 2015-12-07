package routers

import (
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/usersvc/controllers/mq"
	"github.com/byrnedo/apibase/controllers"
)

func InitMq(natsCon *natsio.Nats) {
	controllers.SubscribeNatsRoutes(natsCon, "user_svc_worker", mq.NewHealthcheckController(natsCon))
	controllers.SubscribeNatsRoutes(natsCon, "user_svc_worker", mq.NewUsersController(natsCon))
}

