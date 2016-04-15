package authorizr

import (
	"net/http"
	"github.com/tecsisa/authorizr/api"
)

// Core is the manager of authorizR. This use abstractions of connectors for backends,
// that you define at startup
type Core struct {
	userapi *api.UsersAPI
	groupapi *api.GroupAPI
	policyapi *api.PolicyAPI
}

func NewCore() (*Core){
	return &Core{
		userapi: &api.UsersAPI{},
	}
}

func (core *Core) GetUserAPI() (*api.UsersAPI) {
	return core.userapi
}

func (core *Core) GetGroupAPI() (*api.GroupAPI) {
	return core.groupapi
}

func (core *Core) GetPolicyAPI() (*api.PolicyAPI) {
	return core.policyapi
}

func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
}