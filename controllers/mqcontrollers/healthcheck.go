package mqcontrollers

import (
	"github.com/byrnedo/apibase/natsio"
	r "github.com/byrnedo/apibase/routes"
	"github.com/nats-io/nats"
)

type HealthcheckController struct {
	routes  []*r.NatsRoute
	natsCon *natsio.Nats
}

func (c *HealthcheckController) GetRoutes() []*r.NatsRoute {
	return c.routes
}

func NewHealthcheckController(nc *natsio.Nats) (hc *HealthcheckController) {
	hc = &HealthcheckController{}
	hc.natsCon = nc
	hc.routes = []*r.NatsRoute{
		r.NewNatsRoute("user.healthcheck", hc.Healthcheck),
	}
	return
}

func (c *HealthcheckController) Healthcheck(m *nats.Msg) {
	//c.natsCon.Publish(m.Reply, "up up up")
}
