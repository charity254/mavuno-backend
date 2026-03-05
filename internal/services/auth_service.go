package services

import (
	"fmt"
	"time"

	 "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/mavuno/mavuno-backend/internal/models"
    "github.com/mavuno/mavuno-backend/internal/storage"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *storage.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *storage.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(email, password, fullName, role string) error {
	if role != "farmer" && role != "buyer" {
		return fmt.Errorf("role must be either farmer or buyer")
	}

	if len(password) < 8 {
		return fmt.Errorf("Password must be at least 8 characters")
	}

	if email == "" {
		return fmt.Errorf("Email is required")
	}
	if fullName == "" {
		return fmt.Errorf("Full name is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	user := &models.User{
		ID:           uuid.New(),        // generate a unique ID
        Email:        email,
        PasswordHash: string(hashedPassword),
        Role:         role,
        FullName:     fullName,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        IsActive:     true,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		if isDuplicateEmailError(err) {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role": user.Role,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}

func isDuplicateEmailError(err error) bool {
	return err != nil && len(err.Error()) > 0 &&
	contains(err.Error(), "duplicate key") ||
	contains(err.Error(), "unique constraint")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
	len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0, i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}