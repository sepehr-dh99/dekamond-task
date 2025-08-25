package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dekamond-task/service"
)

type UserController struct {
	userSvc *service.UserService
}

func NewUserController(u *service.UserService) *UserController {
	return &UserController{userSvc: u}
}

// ListUsersHandler handles GET /users.
// @Summary List users
// @Description List users with optional search, pagination (requires JWT).
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param search query string false "Search by phone substring"
// @Success 200 {object} map[string]interface{}
// @Router /users [get]
// @Security     BearerAuth
func (uc *UserController) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(q.Get("size"))
	if size < 1 {
		size = 10
	}
	search := q.Get("search")

	users, total := uc.userSvc.ListUsers(search, page, size)
	// Prepare response
	resp := map[string]interface{}{
		"total": total, "page": page, "size": size, "users": users,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUserHandler handles GET /users/{phone}.
// @Summary Get user by phone
// @Param phone path string true "User phone"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string
// @Router /users/{phone} [get]
// @Security     BearerAuth
func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract phone from URL path, e.g. /users/0912...
	phone := r.URL.Path[len("/users/"):]
	user, exists := uc.userSvc.GetUser(phone)
	if !exists {
		http.Error(w, `{"error":"user not found"}`, 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
