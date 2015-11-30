package mq

import (
	r "github.com/byrnedo/apibase/routes"
	"github.com/apcera/nats"
)

type UsersController struct {
	routes []*r.NatsRoute
	encCon *nats.EncodedConn
}

func (c *UsersController) GetRoutes() []*r.NatsRoute {
	return []*r.NatsRoute{
		r.NewNatsRoute("user.users.get", c.List),
		r.NewNatsRoute("user.users.create", c.List),
		r.NewNatsRoute("user.users.update", c.List),
		r.NewNatsRoute("user.users.delete", c.List),
	}
}

func NewUsersController(nc *nats.EncodedConn) (pC *UsersController) {
	pC = &UsersController{}
	pC.encCon = nc
	return
}

func (c *UsersController) List(m *nats.Msg) {
	c.encCon.Publish(m.Reply, "Not implemented")
}
