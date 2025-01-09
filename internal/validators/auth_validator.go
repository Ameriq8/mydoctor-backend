package validators

import (
	"server/internal/models"
	"server/pkg/utils"

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
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" binding:"required,min=8"`
}

// TLoginRequest represents the structure for login requests
type TLoginRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,phone"`
	Password    string `json:"password" binding:"required"`
}

// ValidateAdapterUser validates a user request
func ValidateAdapterUser(c *gin.Context, user models.User) {
	utils.ValidateRequest(c, user)
}

// ValidateAdapterSession validates a session request
func ValidateAdapterSession(c *gin.Context, session TSession) {
	utils.ValidateRequest(c, session)
}

// ValidateVerificationToken validates a verification token request
func ValidateVerificationToken(c *gin.Context, token TVerificationToken) {
	utils.ValidateRequest(c, token)
}
