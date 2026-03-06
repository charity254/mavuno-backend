package models

import (
	//"time"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			 uuid.UUID `json:"id"`
	Email 		 string	   `json:"email"`
	PasswordHash string	   `json:"-"`
	Role		 string    `json:"role"`
	FullName	 string    `json:"full_name"`
	PhoneNumber  string	   `json:"phone_number"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	IsActive     bool   `json:"is_active"` 
}