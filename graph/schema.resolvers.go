package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/gpiecyk/data-warehouse/graph/generated"
	"github.com/gpiecyk/data-warehouse/graph/model"
	"github.com/gpiecyk/data-warehouse/internal/users"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := &users.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Mobile:    input.Mobile,
		Email:     input.Email,
		Password:  input.Password,
	}

	user, err := r.UserService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Mobile:    user.Mobile,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	user, err := r.UserService.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Mobile:    user.Mobile,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
	}, nil
}

func (r *queryResolver) Users(ctx context.Context, limit *int) ([]*model.User, error) {
	users, err := r.UserService.FindUsersWithLimit(ctx, *limit)
	if err != nil {
		return nil, err
	}

	var resultUsers []*model.User
	for _, user := range *users {
		grahpqlUser := &model.User{
			ID:        int(user.ID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Mobile:    user.Mobile,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt.Time,
		}
		resultUsers = append(resultUsers, grahpqlUser)
	}
	return resultUsers, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
