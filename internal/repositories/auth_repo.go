package repositories

import (
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	// User operations
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id int) (*models.User, error)
	GetUserByEmailOrPhone(emailOrPhone string) (*models.User, error) // Updated method
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id int) error

	// Session operations
	CreateSession(session *models.Session) (*models.Session, error)
	GetSessionAndUser(sessionToken string) (*models.Session, *models.User, error)
	DeleteSession(sessionToken string) error

	// Verification token operations
	CreateVerificationToken(token *models.VerificationToken) (*models.VerificationToken, error)
	UseVerificationToken(identifier, token string) (*models.VerificationToken, error)
}

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

// User operations implementations
func (r *authRepository) CreateUser(user *models.User) (*models.User, error) {
	start := time.Now()

	query := `
		INSERT INTO users (name, email, phone_number, image, email_verified)
		VALUES (:name, :email, :phone_number, :image, :email_verified)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		trackMetrics("CreateUser", "users", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(user); err != nil {
			trackMetrics("CreateUser", "users", start, err)
			return nil, err
		}
	}

	trackMetrics("CreateUser", "users", start, nil)
	return user, nil
}

func (r *authRepository) GetUser(id int) (*models.User, error) {
	start := time.Now()

	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)

	trackMetrics("GetUser", "users", start, err)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByEmailOrPhone(emailOrPhone string) (*models.User, error) {
	start := time.Now()

	var user models.User
	query := `
		SELECT * FROM users 
		WHERE email = $1 OR phone_number = $1`
	err := r.db.Get(&user, query, emailOrPhone)

	trackMetrics("GetUserByEmailOrPhone", "users", start, err)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(user *models.User) (*models.User, error) {
	start := time.Now()

	query := `
		UPDATE users
		SET name = :name, email = :email, phone_number = :phone_number, image = :image, updated_at = NOW()
		WHERE id = :id
		RETURNING *`

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		trackMetrics("UpdateUser", "users", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(user); err != nil {
			trackMetrics("UpdateUser", "users", start, err)
			return nil, err
		}
	}

	trackMetrics("UpdateUser", "users", start, nil)
	return user, nil
}

func (r *authRepository) DeleteUser(id int) error {
	start := time.Now()

	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)

	trackMetrics("DeleteUser", "users", start, err)
	return err
}

// Session operations implementations
func (r *authRepository) CreateSession(session *models.Session) (*models.Session, error) {
	start := time.Now()

	query := `
		INSERT INTO sessions (user_id, expires, session_token)
		VALUES (:user_id, :expires, :session_token)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, session)
	if err != nil {
		trackMetrics("CreateSession", "sessions", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(session); err != nil {
			trackMetrics("CreateSession", "sessions", start, err)
			return nil, err
		}
	}

	trackMetrics("CreateSession", "sessions", start, nil)
	return session, nil
}

func (r *authRepository) GetSessionAndUser(sessionToken string) (*models.Session, *models.User, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("GetSessionAndUser", "sessions", start, err)
		return nil, nil, err
	}
	defer tx.Rollback()

	var session models.Session
	sessionQuery := `SELECT * FROM sessions WHERE session_token = $1`
	err = tx.Get(&session, sessionQuery, sessionToken)
	if err != nil {
		trackMetrics("GetSessionAndUser", "sessions", start, err)
		return nil, nil, err
	}

	var user models.User
	userQuery := `SELECT * FROM users WHERE id = $1`
	err = tx.Get(&user, userQuery, session.UserID)
	if err != nil {
		trackMetrics("GetSessionAndUser", "sessions", start, err)
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		trackMetrics("GetSessionAndUser", "sessions", start, err)
		return nil, nil, err
	}

	trackMetrics("GetSessionAndUser", "sessions", start, nil)
	return &session, &user, nil
}

func (r *authRepository) DeleteSession(sessionToken string) error {
	start := time.Now()

	query := `DELETE FROM sessions WHERE session_token = $1`
	_, err := r.db.Exec(query, sessionToken)

	trackMetrics("DeleteSession", "sessions", start, err)
	return err
}

// Verification token operations implementations
func (r *authRepository) CreateVerificationToken(token *models.VerificationToken) (*models.VerificationToken, error) {
	start := time.Now()

	query := `
		INSERT INTO verification_tokens (id, token, expires)
		VALUES (:id, :token, :expires)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, token)
	if err != nil {
		trackMetrics("CreateVerificationToken", "verification_tokens", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(token); err != nil {
			trackMetrics("CreateVerificationToken", "verification_tokens", start, err)
			return nil, err
		}
	}

	trackMetrics("CreateVerificationToken", "verification_tokens", start, nil)
	return token, nil
}

func (r *authRepository) UseVerificationToken(identifier, token string) (*models.VerificationToken, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("UseVerificationToken", "verification_tokens", start, err)
		return nil, err
	}
	defer tx.Rollback()

	var verificationToken models.VerificationToken
	query := `SELECT * FROM verification_tokens WHERE identifier = $1 AND token = $2`
	err = tx.Get(&verificationToken, query, identifier, token)
	if err != nil {
		trackMetrics("UseVerificationToken", "verification_tokens", start, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		trackMetrics("UseVerificationToken", "verification_tokens", start, err)
		return nil, err
	}

	trackMetrics("UseVerificationToken", "verification_tokens", start, nil)
	return &verificationToken, nil
}
