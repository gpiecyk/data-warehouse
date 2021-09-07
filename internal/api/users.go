package api

import (
	"context"

	"github.com/gpiecyk/data-warehouse/internal/users"
)

func (api *API) CreateUser(ctx context.Context, user *users.User) (*users.User, error) {
	return api.userService.CreateUser(ctx, user)
}

func (api *API) UpdateUser(ctx context.Context, user *users.User, id int) (*users.User, error) {
	return api.userService.UpdateUser(ctx, user, id)
}

func (api *API) DeleteUser(ctx context.Context, id int) error {
	return api.userService.DeleteUser(ctx, id)
}

func (api *API) GetUserById(ctx context.Context, id int) (*users.User, error) {
	return api.userService.GetUserById(ctx, id)
}

func (api *API) FindUsersWithLimit(ctx context.Context, limit int) (*[]users.User, error) {
	return api.userService.FindUsersWithLimit(ctx, limit)
}
