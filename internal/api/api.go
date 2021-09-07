package api

import (
	"time"

	"github.com/gpiecyk/data-warehouse/internal/users"
)

var (
	now = time.Now()
)

type API struct {
	userService *users.UserService
}

func (a *API) Health() (map[string]interface{}, error) {
	return map[string]interface{}{
		"env":        "testing",
		"version":    "v0.1.0",
		"commit":     "<git commit hash>",
		"status":     "all systems up and running",
		"startedAt":  now.String(),
		"releasedOn": now.String(),
	}, nil
}

func NewService(us *users.UserService) (*API, error) {
	return &API{
		userService: us,
	}, nil
}
