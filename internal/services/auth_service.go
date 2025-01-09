package services

import (
	"errors"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/validators"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go" // JWT library
	"golang.org/x/crypto/bcrypt"  // to hash and compare passwords
)

type AuthService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// User operations
func (s *AuthService) CreateUser(user *validators.TRegisterRequest) (*models.User, error) {
	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	modelUser := mapUserToModel(user)
	modelUser.Password = string(hashedPassword)

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmailOrPhone(user.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	createdUser, err := s.repo.CreateUser(modelUser)
	if err != nil {
		return nil, err
	}

	return mapModelToUser(createdUser), nil
}

// AuthenticateUser validates the user's credentials (login)
func (s *AuthService) AuthenticateUser(loginRequest *validators.TLoginRequest) (*models.User, error) {
	user, err := s.repo.GetUserByEmailOrPhone(loginRequest.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	// Convert string userID to integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Fetch the user from the repository using the user ID
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GenerateAuthToken generates a JWT token for the authenticated user
func (s *AuthService) GenerateAuthToken(user *models.User) (string, error) {
	// Define claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte("your_secret_key") // Replace with a secret key from config

	// Generate signed token string
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Verification token operations (unchanged)
func (s *AuthService) CreateVerificationToken(token *validators.TVerificationToken) (*validators.TVerificationToken, error) {
	modelToken, err := mapVerificationTokenToModel(token)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.CreateVerificationToken(modelToken)
	if err != nil {
		return nil, err
	}

	return mapModelToVerificationToken(created), nil
}

// VerifyToken verifies a token (e.g., email/phone verification)
func (s *AuthService) VerifyToken(identifier, token string) (bool, error) {
	// Fetch the verification token from the repository using the identifier and token
	verificationToken, err := s.repo.UseVerificationToken(identifier, token)
	if err != nil {
		return false, err
	}
	if verificationToken == nil {
		return false, errors.New("invalid token")
	}

	// Check if the token has expired
	if verificationToken.Expires.Before(time.Now()) {
		return false, errors.New("token has expired")
	}

	return true, nil
}

// Mapping functions

// Map user to model
func mapUserToModel(user *validators.TRegisterRequest) *models.User {
	return &models.User{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: &user.PhoneNumber, // Ensure this field exists in TRegisterRequest
		Password:    user.Password,
	}
}

// Map model user to validator user
func mapModelToUser(user *models.User) *models.User {
	return &models.User{
		BaseModel:     user.BaseModel, // Include the embedded BaseModel
		Name:          user.Name,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		EmailVerified: user.EmailVerified,
		Image:         user.Image,
	}
}

// Map verification token to model
func mapVerificationTokenToModel(token *validators.TVerificationToken) (*models.VerificationToken, error) {
	expires, err := time.Parse(time.RFC3339, token.Expires)
	if err != nil {
		return nil, err
	}

	return &models.VerificationToken{
		Token:   token.Token,
		Expires: expires,
	}, nil
}

// Map model verification token to validator verification token
func mapModelToVerificationToken(token *models.VerificationToken) *validators.TVerificationToken {
	return &validators.TVerificationToken{
		ID:      token.ID, // Use the ID from BaseModel
		Token:   token.Token,
		Expires: token.Expires.Format(time.RFC3339),
	}
}
