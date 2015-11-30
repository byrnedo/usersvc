package routers

import (
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/usersvc/controllers/mq"
	"github.com/byrnedo/apibase/controllers"
)

func InitMq(natsCon *natsio.Nats) {
	controllers.SubscribeNatsRoutes(natsCon, mq.NewHealthcheckController(natsCon.EncCon))
	controllers.SubscribeNatsRoutes(natsCon, mq.NewUsersController(natsCon.EncCon))
}
