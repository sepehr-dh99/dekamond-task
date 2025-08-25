package response

import (
	"encoding/json"
	"net/http"
)

// Generic API response wrapper
type Response[T any] struct {
	Success bool   `json:"success" example:"true"`
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message" example:"OK"`
}

// ErrorResponse is used for non-200 responses (no `data`)
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Invalid request"`
}

// JSON sends a JSON response with status code.
func JSON[T any](w http.ResponseWriter, status int, success bool, data *T, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response[T]{
		Success: success,
		Data:    data,
		Message: message,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// Success helper for success responses.
func Success[T any](w http.ResponseWriter, data *T, message string) {
	JSON(w, http.StatusOK, true, data, message)
}

// Error helper for error responses.
func Error(w http.ResponseWriter, status int, message string) {
	JSON[any](w, status, false, nil, message)
}

type PaginatedResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Total   int    `json:"total"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Data    []T    `json:"data"`
}

func Paginated[T any](w http.ResponseWriter, data []T, total, page, size int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := PaginatedResponse[T]{
		Success: true,
		Message: message,
		Total:   total,
		Page:    page,
		Size:    size,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}
