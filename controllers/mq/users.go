package mq

import (
	r "github.com/byrnedo/apibase/routes"
	"github.com/apcera/nats"
	"github.com/byrnedo/usersvc/msgspec/mq"
)

type UsersController struct {
	routes []*r.NatsRoute
	encCon *nats.EncodedConn
}


func (c *UsersController) GetRoutes() []*r.NatsRoute {
	return []*r.NatsRoute{
		r.NewNatsRoute(mq.AuthenticateUserSubject, c.Authenticate),
	}
}

func NewUsersController(nc *nats.EncodedConn) (pC *UsersController) {
	pC = &UsersController{}
	pC.encCon = nc
	return
}

func (c *UsersController) Authenticate(subj string, reply string, data *mq.AuthenticateUserRequest) {
	c.encCon.Publish(reply, "Not implemented")
}
