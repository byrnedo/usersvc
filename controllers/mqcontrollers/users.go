package mqcontrollers

import (
	"github.com/byrnedo/apibase/db/mongo/defaultmongo"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
	r "github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/svccommon/validate"
	"github.com/byrnedo/usersvc/daos"
	"github.com/byrnedo/usersvc/msgspec/mqmsgspec"
	"github.com/nats-io/nats"
)

type UsersController struct {
	routes    []*r.NatsRoute
	natsCon   *natsio.Nats
	userModel daos.UserDAO
}

func (c *UsersController) GetRoutes() []*r.NatsRoute {
	return []*r.NatsRoute{
		r.NewNatsRoute(mqmsgspec.AuthenticateUserSubject, c.Authenticate),
	}
}

func NewUsersController(nc *natsio.Nats, encryptionKey string) (pC *UsersController) {
	pC = &UsersController{}
	pC.natsCon = nc
	pC.userModel = daos.NewDefaultUserDAO(defaultmongo.Conn(), encryptionKey)
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

func (c *UsersController) Authenticate(subj string, reply string, data *mqmsgspec.InnerAuthenticateUserRequest) {
	Info.Println("Got authenticate request:", data)

	var valid bool
	if valErrs := validate.ValidateStruct(data); len(valErrs) != 0 {
		for key, fieldErr := range valErrs {
			Error.Println("Validation failed:", key, ":", fieldErr.Tag)

		}
		valid = false
	} else {
		err := c.userModel.Authenticate(data.GetUsername(), data.GetPassword())

		if err == nil {
			valid = true
			Info.Println("Authentication successful")
		} else {
			valid = false
			Info.Println("Authentication failed:", err)
		}
	}

	response := mqmsgspec.NewAuthenticateUserResponse(&mqmsgspec.InnerAuthenticateUserResponse{Authenticated: &valid})
	if err := c.natsCon.Publish(reply, data.GetContext(), response); err != nil {
		Error.Println("Error sending reply:" + err.Error())
	}
}
