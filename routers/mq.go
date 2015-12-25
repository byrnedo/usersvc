package routers

import (
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/natsio/defaultnats"
	"github.com/byrnedo/usersvc/controllers/mqcontrollers"
)

func init() {

	encryptionKey, err := config.Conf.GetString("encryption-key")
	if err != nil {
		panic("Failed to get encryption-key:" + err.Error())
	}

	controllers.SubscribeNatsRoutes(defaultnats.Conn, "user_svc_worker", mqcontrollers.NewHealthcheckController(defaultnats.Conn))
	controllers.SubscribeNatsRoutes(defaultnats.Conn, "user_svc_worker", mqcontrollers.NewUsersController(defaultnats.Conn, encryptionKey))
}
