// Package api provides HTTP API layer
package api

import (
	"awdoc/examples/complex/service"
	"fmt"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler creates a new handler
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// HandleGetUser handles GET /users/:id request
func (h *UserHandler) HandleGetUser(id int) (string, error) {
	user, err := h.svc.GetUser(id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("User: %s (%s)", user.Name, user.Email), nil
}

// HandleCreateUser handles POST /users request
func (h *UserHandler) HandleCreateUser(name, email string) error {
	return h.svc.RegisterUser(name, email)
}

// HandleListUsers handles GET /users request
func (h *UserHandler) HandleListUsers() ([]string, error) {
	users, err := h.svc.ListUsers()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, user := range users {
		result = append(result, fmt.Sprintf("%d: %s", user.ID, user.Name))
	}
	return result, nil
}

// HandleDeleteUser handles DELETE /users/:id request
func (h *UserHandler) HandleDeleteUser(id int) error {
	return h.svc.DeleteUser(id)
}
