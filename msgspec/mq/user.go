package mq
import "github.com/byrnedo/apibase/natsio"

const (
	AuthenticateUserSubject = "user.users.authenticate"
)

type AuthenticateUserRequest struct {
	natsio.NatsDTO
	User string
	Password string
}
