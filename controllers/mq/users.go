package mq

import (
	r "github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/usersvc/msgspec/mq"
	"github.com/byrnedo/usersvc/models"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
"github.com/apcera/nats"
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

func (c *UsersController) Find(m *nats.Msg) {
}

func (c *UsersController) Create(m *nats.Msg) {
}

func (c *UsersController) Update(m *nats.Msg) {
}

func (c *UsersController) Delete(m *nats.Msg) {
}

func (c *UsersController) Authenticate(subj string, reply string, data *mq.AuthenticateUserRequest) {
	valid:= c.userModel.Authenticate(data.User, data.Password)
	response := mq.AuthenticateUserResponse{
		NatsDTO: natsio.NatsDTO{NatsCtx: data.NatsCtx},
		Authenticated: valid,
	}

	if err:=c.natsCon.Publish(reply, &response); err != nil {
		Error.Println("Error sending reply:" + err.Error())
	}
}
