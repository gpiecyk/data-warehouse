package users

import (
	"context"

	"gorm.io/gorm"
)

// UserDao
// jeśli mam wiele encji w users to powinieniem stworzyć folder "repository" i w nim umieścić user.go, address.go, etc.

type repository interface { // zmienic nazwe, bo gryzie sie z userRepository!!!!!!!!!! Moze datastore?
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User, id int) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*User, error)
	FindWithLimit(ctx context.Context, limit int) (*[]User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (repository *userRepository) Create(ctx context.Context, user *User) error {
	result := repository.db.Create(user)
	return result.Error
}

func (repository *userRepository) Update(ctx context.Context, user *User, id int) error {
	result := repository.db.Save(user)
	return result.Error
}

func (repository *userRepository) Delete(ctx context.Context, id int) error {
	result := repository.db.Delete(&User{}, id)
	return result.Error
}

func (repository *userRepository) GetById(ctx context.Context, id int) (*User, error) {
	user := User{}
	result := repository.db.First(&user, id)
	return &user, result.Error
}

func (repository *userRepository) FindWithLimit(ctx context.Context, limit int) (*[]User, error) {
	users := new([]User)
	result := repository.db.Order("ID ASC").Limit(limit).Find(users)
	return users, result.Error
}

func newRepository(db *gorm.DB) (repository, error) {
	return &userRepository{
		db: db,
	}, nil
}
