package handlers

import (
	"net/http"
	"server/internal/services"
	"server/internal/validators"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.RegisterUser)                                         // Register a new user
	r.POST("/login", h.LoginUser)                                               // User login
	r.GET("/me", h.GetAuthenticatedUser)                                        // Get authenticated user's details
	r.DELETE("/logout", h.LogoutUser)                                           // Logout user
	r.POST("/verification-tokens", h.CreateVerificationToken)                   // Create a verification token
	r.DELETE("/verification-tokens/:identifier/:token", h.UseVerificationToken) // Verify a token
}

// RegisterUser handles user registration requests
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req validators.TRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.service.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// LoginUser handles user login requests
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req validators.TLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Attempt to authenticate the user
	user, err := h.service.LoginUser(&req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Generate a session or JWT token for the user
	token, err := h.service.GenerateAuthToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	// Respond with the user information and the generated token
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// GetAuthenticatedUser fetches the currently authenticated user's data
func (h *AuthHandler) GetAuthenticatedUser(c *gin.Context) {
	// In a real-world scenario, you would get the user ID from the token in the request header
	// For now, we will mock the user retrieval
	userID := c.GetString("user_id") // Assuming you've set the user ID in the context earlier
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch the user from the database using the service (or session)
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// LogoutUser handles user logout and session destruction
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	// Logic to invalidate the token or session
	// Assuming the token is stored in a cookie or header and needs to be invalidated
	// Invalidate session or JWT token
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// CreateVerificationToken creates a token for email/phone verification
func (h *AuthHandler) CreateVerificationToken(c *gin.Context) {
	var req validators.TVerificationToken
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token, err := h.service.CreateVerificationToken(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// UseVerificationToken verifies and uses a verification token
func (h *AuthHandler) UseVerificationToken(c *gin.Context) {
	identifier := c.Param("identifier")
	token := c.Param("token")

	// Logic to verify the token
	isVerified, err := h.service.VerifyToken(identifier, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token verification failed"})
		return
	}

	if !isVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token verified successfully"})
}
