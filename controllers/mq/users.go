package mq

import (
	r "github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/usersvc/msgspec/mq"
	"github.com/byrnedo/usersvc/models"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
)

type UsersController struct {
	routes    []*r.NatsRoute
	natsCon   *natsio.Nats
	userModel models.UserModel
}


func (c *UsersController) GetRoutes() []*r.NatsRoute {
	return []*r.NatsRoute{
		r.NewNatsRoute(mq.AuthenticateUserSubject, c.Authenticate),
	}
}

func NewUsersController(nc *natsio.Nats) (pC *UsersController) {
	pC = &UsersController{}
	pC.natsCon = nc
	pC.userModel = models.NewDefaultUserModel()
	return
}

func (c *UsersController) Authenticate(subj string, reply string, data *mq.AuthenticateUserRequest) {
	valid, err := c.userModel.Authenticate(data.User, data.Password)
	if err != nil {
		Error.Println("Error on login:" + err.Error())
	}

	response := mq.AuthenticateUserResponse{
		Authenticated: valid,
	}

	if err:=c.natsCon.Publish(reply, &response); err != nil {
		Error.Println("Error sending reply:" + err.Error())
	}
}
