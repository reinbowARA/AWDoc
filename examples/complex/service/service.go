// Package service provides business logic layer
package service

import (
	"awdoc/examples/complex/repository"
)

// UserService provides user business logic
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(name, email string) error {
	user := &repository.User{
		Name:  name,
		Email: email,
	}
	return s.repo.Save(user)
}

// GetUser fetches user information
func (s *UserService) GetUser(id int) (*repository.User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id int, name, email string) error {
	user := &repository.User{
		ID:    id,
		Name:  name,
		Email: email,
	}
	return s.repo.Save(user)
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// ListUsers returns all users
func (s *UserService) ListUsers() ([]*repository.User, error) {
	return s.repo.GetAll()
}
