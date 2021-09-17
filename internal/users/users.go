package users

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gpiecyk/data-warehouse/internal/platform/cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

// wrzucic do pliku service.go, poniewaz w service nie powinienem odnosić się do bazy danych (gorm)
type UserService struct {
	repository repository
	cache      *cache.Client
}

// handler_users.go -> api user handler - first contact
// api.users -> api + service (calls business logic in users.go)
// users.go (this file) -> business logic (maybe add that stuff into service.go)
// repository.go -> dao
func (service *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	err = service.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) UpdateUser(ctx context.Context, user *User, id int) (*User, error) {
	if user.ID == 0 || user.ID != uint(id) {
		return nil, errors.New("missing/wrong 'id' in the body")
	}

	if _, err := service.GetUserById(ctx, id); err != nil {
		return nil, err
	}

	if err := service.repository.Update(ctx, user, id); err != nil {
		return nil, err
	}

	// delete is an idempotent operation - always a safe action
	// it does not impair consistency
	cacheErr := service.cache.Delete(ctx, generateUserKeyWithId(id))
	if cacheErr != nil {
		log.Println(cacheErr)
	}
	return user, nil
}

func (service *UserService) DeleteUser(ctx context.Context, id int) error {
	dbErr := service.repository.Delete(ctx, id)
	if err := service.cache.Delete(ctx, generateUserKeyWithId(id)); err != nil {
		log.Println(err)
	}
	return dbErr
}

func (service *UserService) GetUserById(ctx context.Context, id int) (*User, error) {
	user := new(User)
	key := generateUserKeyWithId(id)
	if err := service.cache.Get(ctx, key, &user); err != nil {
		log.Println(err, "| method: GetUserById | id:", id)
	} else {
		return user, nil
	}

	user, err := service.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	cacheItem := &cache.Item{Key: key, Value: user}
	if err := service.cache.Set(ctx, cacheItem); err != nil {
		log.Printf("cache set error: %v", err)
	}

	return user, nil
}

func (service *UserService) FindUsersWithLimit(ctx context.Context, limit int) (*[]User, error) {
	users, err := service.repository.FindWithLimit(ctx, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Authenticate returns user's id if provided password and existing password match
func (service *UserService) Authenticate(ctx context.Context, email, password string) (int, error) {
	user, err := service.repository.GetByEmail(ctx, email)
	if err != nil {
		log.Println("authentication failed, error", err)
		return -1, err
	}

	passwordsMatchErr := CompareHashedPasswordWithPlainTextPassword(user.Password, password)
	return int(user.ID), passwordsMatchErr
}

func generateUserKeyWithId(id int) string {
	return fmt.Sprintf("user:id:%v", id)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashedPasswordWithPlainTextPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewService(db *gorm.DB, cache *cache.Client) (*UserService, error) {
	repositoryObj, err := newRepository(db)
	if err != nil {
		return nil, err
	}

	return &UserService{
		repository: repositoryObj,
		cache:      cache,
	}, nil
}
