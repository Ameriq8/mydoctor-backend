package validators

import (
	"server/internal/models"
	"server/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

// AdapterSession represents the structure for session operations
type TSession struct {
	SessionToken string `json:"sessionToken" validate:"required"`
	UserID       string `json:"userId" validate:"required"`
	Expires      string `json:"expires" validate:"required"`
}

// VerificationToken represents the structure for verification token operations
type TVerificationToken struct {
	ID      int64  `json:"id" validate:"required"`
	Token   string `json:"token" validate:"required"`
	Expires string `json:"expires" validate:"required"`
}

// TRegisterRequest represents the structure for registration requests
type TRegisterRequest struct {
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,phone"`
	Password    string `json:"password" binding:"required"`
}

// TLoginRequest represents the structure for login requests
type TLoginRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,phone"`
	Password    string `json:"password" binding:"required"`
}

// ValidateAdapterUser validates a user request
func ValidateAdapterUser(c *gin.Context, user models.User) {
	middlewares.ValidateRequest(c, user)
}

// ValidateAdapterSession validates a session request
func ValidateAdapterSession(c *gin.Context, session TSession) {
	middlewares.ValidateRequest(c, session)
}

// ValidateVerificationToken validates a verification token request
func ValidateVerificationToken(c *gin.Context, token TVerificationToken) {
	middlewares.ValidateRequest(c, token)
}
