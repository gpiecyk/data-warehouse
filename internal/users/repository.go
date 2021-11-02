package users

import (
	"context"

	"gorm.io/gorm"
)

type repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User, id int) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	FindWithLimit(ctx context.Context, limit int) (*[]User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (repository *userRepository) Create(ctx context.Context, user *User) error {
	result := repository.db.WithContext(ctx).Create(user)
	return result.Error
}

func (repository *userRepository) Update(ctx context.Context, user *User, id int) error {
	result := repository.db.WithContext(ctx).Save(user)
	return result.Error
}

func (repository *userRepository) Delete(ctx context.Context, id int) error {
	result := repository.db.WithContext(ctx).Delete(&User{}, id)
	return result.Error
}

func (repository *userRepository) GetById(ctx context.Context, id int) (*User, error) {
	user := User{}
	result := repository.db.WithContext(ctx).First(&user, id)
	return &user, result.Error
}

func (repository *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}
	result := repository.db.WithContext(ctx).Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (repository *userRepository) FindWithLimit(ctx context.Context, limit int) (*[]User, error) {
	users := new([]User)
	result := repository.db.WithContext(ctx).Order("ID ASC").Limit(limit).Find(users)
	return users, result.Error
}

func newRepository(db *gorm.DB) (repository, error) {
	return &userRepository{
		db: db,
	}, nil
}
