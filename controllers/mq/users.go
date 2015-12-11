package mq

import (
	"github.com/nats-io/nats"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
	r "github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/usersvc/models"
	"github.com/byrnedo/usersvc/msgspec/mq"
	"github.com/byrnedo/apibase/natsio/protobuf"
)

type UsersController struct {
	routes    []*r.NatsRoute
	natsCon   *natsio.Nats
	userModel models.UserModel
}

type WrapAuthUserRes struct {
	*mq.AuthenticateUserResponse
}
func (w *WrapAuthUserRes) SetContext(ctx *protobuf.NatsContext) {
	w.Context = ctx
}
func newWrapAuthUserRes(msg *mq.AuthenticateUserResponse) *WrapAuthUserRes {
	return &WrapAuthUserRes{msg}
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
	valid := c.userModel.Authenticate(data.GetUsername(), data.GetPassword())
	response := mq.AuthenticateUserResponse{
		Authenticated: &valid,
	}

	if err := c.natsCon.Publish(reply, data.GetContext(), newWrapAuthUserRes(&response)); err != nil {
		Error.Println("Error sending reply:" + err.Error())
	}
}
