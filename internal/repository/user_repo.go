package repository

import (
	"errors"
	"test-project/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type InMemoryUserRepo struct {
	data map[string]domain.User
}

func NewInMemoryUserRepo() domain.UserRepository {
	return &InMemoryUserRepo{
		data: make(map[string]domain.User),
	}
}

func (r *InMemoryUserRepo) FindAll() ([]domain.User, error) {
	users := make([]domain.User, 0, len(r.data))
	for _, u := range r.data {
		users = append(users, u)
	}
	return users, nil
}

func (r *InMemoryUserRepo) FindByID(id string) (domain.User, error) {
	u, ok := r.data[id]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return u, nil
}

func (r *InMemoryUserRepo) Create(u domain.User) (domain.User, error) {
	r.data[u.ID] = u
	return u, nil
}
