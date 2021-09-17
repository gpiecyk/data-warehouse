package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/gpiecyk/data-warehouse/graph/generated"
	"github.com/gpiecyk/data-warehouse/graph/model"
	"github.com/gpiecyk/data-warehouse/internal/auth"
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

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	id, err := r.UserService.Authenticate(ctx, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	id, err := auth.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}

	token, err := auth.GenerateToken(id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	if user := auth.GetUserFromContext(ctx); user == nil {
		return &model.User{}, fmt.Errorf("access denied")
	}

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
	if user := auth.GetUserFromContext(ctx); user == nil {
		return nil, fmt.Errorf("access denied")
	}

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
