package models

import (
	//"time"
	"github.com/google/uuid"
)

type User struct {
	ID 			 uuid.UUID `json:"id"`
	Email 		 string	   `json:"email"`
	PasswordHash string	   `json:"-"`
	Role		 string    `json:"role"`
	FullName	 string    `json:"full_name"`
	PhoneNumber  string	   `json:"phone_number"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	IsActive     string    `json:"is_active"` 
}