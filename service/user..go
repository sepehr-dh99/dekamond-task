package service

import (
	"strings"
	"sync"
	"time"

	"dekamond-task/model"
)

type UserService struct {
	mu    sync.RWMutex
	users map[string]model.User // phone -> User
}

func NewUserService() *UserService {
	return &UserService{users: make(map[string]model.User)}
}

// RegisterIfNotExists adds the user if new, and returns the user.
func (u *UserService) RegisterIfNotExists(phone string) model.User {
	u.mu.Lock()
	defer u.mu.Unlock()
	if usr, exists := u.users[phone]; exists {
		return usr
	}
	newUser := model.User{Phone: phone, RegisteredAt: time.Now()}
	u.users[phone] = newUser
	return newUser
}

func (u *UserService) GetUser(phone string) (model.User, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	usr, exists := u.users[phone]
	return usr, exists
}

// ListUsers returns users filtered by search and paginated.
func (u *UserService) ListUsers(search string, page, size int) (result []model.User, total int) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	// Collect all matching users
	for _, usr := range u.users {
		if search == "" || strings.Contains(usr.Phone, search) {
			result = append(result, usr)
		}
	}
	total = len(result)
	// Sort or limit; for simplicity assume insertion order (map is random).
	// Apply offset pagination:
	start := (page - 1) * size
	if start > total {
		return nil, total
	}
	end := start + size
	if end > total {
		end = total
	}
	return result[start:end], total
}
