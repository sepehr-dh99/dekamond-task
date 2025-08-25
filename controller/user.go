package controller

import (
	"net/http"
	"strconv"

	"dekamond-task/controller/dto"
	"dekamond-task/package/response"
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
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param search query string false "Search by phone substring"
// @Success 200 {object} response.PaginatedResponse[dto.UserResponse] "Paginated list of users"
// @Failure 400 {object} response.ErrorResponse "Invalid query"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /users [get]
// @Security BearerAuth
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

	out := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.UserResponse{
			Phone:        u.Phone,
			RegisteredAt: u.RegisteredAt,
		})
	}

	response.Paginated[dto.UserResponse](w, out, total, page, size, "Users fetched successfully")
}

// GetUserHandler handles GET /users/{phone}.
// @Summary Get user by phone
// @Description Retrieve a single user by phone number.
// @Tags Users
// @Accept json
// @Produce json
// @Param phone path string true "User phone"
// @Success 200 {object} response.Response[dto.UserResponse] "User details"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /users/{phone} [get]
// @Security BearerAuth
func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Path[len("/users/"):]
	u, ok := uc.userSvc.GetUser(phone)
	if !ok {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	payload := dto.UserResponse{
		Phone:        u.Phone,
		RegisteredAt: u.RegisteredAt,
	}
	response.Success(w, &payload, "User fetched successfully")
}
