package models

import "time"

// User represents a user in the system
type User struct {
	ID   int32     `json:"id"`
	Name string    `json:"name"`
	DOB  time.Time `json:"dob"`
	Age  int       `json:"age,omitempty"` // Calculated field, not stored in DB
}

// CreateUserRequest is the request body for creating a user
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is the request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse is a generic success response wrapper
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

