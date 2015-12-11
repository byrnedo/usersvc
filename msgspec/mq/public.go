package mq

import "github.com/byrnedo/apibase/natsio/protobuf"

// Wrapper for protobuf generated structs to
// add setters
type AuthenticateUserResponse struct {
	*InnerAuthenticateUserResponse
}

func (w *AuthenticateUserResponse) SetContext(ctx *protobuf.NatsContext) {
	w.Context = ctx
}

func NewAuthenticateUserResponse(r *InnerAuthenticateUserResponse) *AuthenticateUserResponse {
	return &AuthenticateUserResponse{r}
}

type AuthenticateUserRequest struct {
	*InnerAuthenticateUserRequest
}

func (w *AuthenticateUserRequest) SetContext(ctx *protobuf.NatsContext) {
	w.Context = ctx
}

func NewAuthenticateUserRequest(r *InnerAuthenticateUserRequest) *AuthenticateUserRequest {
	return &AuthenticateUserRequest{r}
}
