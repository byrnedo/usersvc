package mq

import (
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
	r "github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/usersvc/models"
	"github.com/byrnedo/usersvc/msgspec/mq"
	"github.com/nats-io/nats"
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

func (c *UsersController) Authenticate(subj string, reply string, data *mq.InnerAuthenticateUserRequest) {
	Info.Println("Got authenticate request:", data)
	err := c.userModel.Authenticate(data.GetUsername(), data.GetPassword())

	var valid bool
	if err == nil {
		valid = true
		Info.Println("Authentication successful")
	} else {
		valid = false
		Info.Println("Authentication failed:", err)
	}
	response := mq.NewAuthenticateUserResponse(&mq.InnerAuthenticateUserResponse{Authenticated: &valid})
	if err := c.natsCon.Publish(reply, data.GetContext(), response); err != nil {
		Error.Println("Error sending reply:" + err.Error())
	}
}
