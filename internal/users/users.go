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
	cacheErr := service.cache.Delete(ctx, generateUserKey(id))
	if cacheErr != nil {
		log.Println(cacheErr)
	}
	return user, nil
}

func (service *UserService) DeleteUser(ctx context.Context, id int) error {
	dbErr := service.repository.Delete(ctx, id)
	if err := service.cache.Delete(ctx, generateUserKey(id)); err != nil {
		log.Println(err)
	}
	return dbErr
}

func (service *UserService) GetUserById(ctx context.Context, id int) (*User, error) {
	user := new(User)
	key := generateUserKey(id)
	if err := service.cache.Get(ctx, key, &user); err != nil {
		log.Println(err, "| method: GetUserById | id:", id)
	} else {
		return user, nil
	}

	user, err := service.repository.GetById(ctx, id)
	if err != nil {
		return nil, &ErrRecordNotFound{
			Message:    fmt.Sprintf("user with id \"%v\" does not exist", id),
			StatusCode: 404, // wrong! not here! this should be in the http handler
			Err:        err,
		}
	}

	cacheItem := &cache.Item{Key: key, Value: user}
	if err := service.cache.Set(ctx, cacheItem); err != nil {
		log.Printf("ERROR cache set error: %v", err)
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

func generateUserKey(id int) string {
	return fmt.Sprintf("user:id:%v", id)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
