package mq
import (
	r "github.com/byrnedo/apibase/routes"
	"github.com/apcera/nats"
)

type HealthcheckController struct {
	routes []*r.NatsRoute
	encCon *nats.EncodedConn
}

func (c *HealthcheckController) GetRoutes() []*r.NatsRoute {
	return c.routes
}

func NewHealthcheckController(nc *nats.EncodedConn) (hc *HealthcheckController) {
	hc = &HealthcheckController{}
	hc.encCon = nc
	hc.routes = []*r.NatsRoute{
		r.NewNatsRoute("user.healthcheck", hc.Healthcheck),
	}
	return
}

func (c *HealthcheckController) Healthcheck(m *nats.Msg) {
	c.encCon.Publish(m.Reply, "up up up")
}



