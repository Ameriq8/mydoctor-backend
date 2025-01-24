package repositories_test

import (
	"server/internal/models"
	"server/internal/repositories"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB

func setupDB(t *testing.T) {
	dsn := "host=localhost port=3002 user=postgres password=Ameriq81 dbname=mydoctor sslmode=disable"

	var err error
	db, err = sqlx.Connect("postgres", dsn)
	require.NoError(t, err)
}

func tearDownDB(t *testing.T) {
	_, err := db.Exec("DELETE FROM users")
	require.NoError(t, err)
}

func setupAndTearDown(t *testing.T) {
	setupDB(t)
	tearDownDB(t)
}

func TestCreateUser(t *testing.T) {
	setupAndTearDown(t)

	repo := repositories.NewAuthRepository(db)

	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Image:    "image.png",
		Password: "tewdhhjcg",
	}

	result, err := repo.CreateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john.doe@example.com", result.Email)
}

func TestGetUser(t *testing.T) {
	setupAndTearDown(t)

	repo := repositories.NewAuthRepository(db)

	user := &models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Image: "image.png",
	}
	createdUser, err := repo.CreateUser(user)
	require.NoError(t, err)

	fetchedUser, err := repo.GetUser(int(createdUser.ID))

	require.NoError(t, err)
	assert.Equal(t, "John Doe", fetchedUser.Name)
	assert.Equal(t, "john.doe@example.com", fetchedUser.Email)
}

func TestUpdateUser(t *testing.T) {
	setupAndTearDown(t)

	repo := repositories.NewAuthRepository(db)

	user := &models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Image: "image.png",
	}
	createdUser, err := repo.CreateUser(user)
	require.NoError(t, err)

	createdUser.Name = "John Updated"
	createdUser.Email = "john.updated@example.com"
	updatedUser, err := repo.UpdateUser(createdUser)

	require.NoError(t, err)
	assert.Equal(t, "John Updated", updatedUser.Name)
	assert.Equal(t, "john.updated@example.com", updatedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	setupAndTearDown(t)

	repo := repositories.NewAuthRepository(db)

	user := &models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Image: "image.png",
	}
	createdUser, err := repo.CreateUser(user)
	require.NoError(t, err)

	err = repo.DeleteUser(int(createdUser.ID))

	require.NoError(t, err)

	_, err = repo.GetUser(int(createdUser.ID))
	require.Error(t, err)
}

func TestGetUserByEmailOrPhone(t *testing.T) {
	setupAndTearDown(t)

	repo := repositories.NewAuthRepository(db)

	user := &models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		PhoneNumber: func() *string {
			s := "1234567890"
			return &s
		}(),
		Image: "image.png",
	}
	_, err := repo.CreateUser(user)
	require.NoError(t, err)

	fetchedUser, err := repo.GetUserByEmailOrPhone("john.doe@example.com")

	require.NoError(t, err)
	assert.Equal(t, "John Doe", fetchedUser.Name)
	assert.Equal(t, "john.doe@example.com", fetchedUser.Email)
}

func TestCleanup(t *testing.T) {
	setupAndTearDown(t)

	repo := repositories.NewAuthRepository(db)

	user := &models.User{
		Name:  "Jane Doe",
		Email: "jane.doe@example.com",
		Image: "image2.png",
	}
	_, err := repo.CreateUser(user)
	require.NoError(t, err)

	tearDownDB(t)

	var users []models.User
	err = db.Select(&users, "SELECT * FROM users")
	require.NoError(t, err)
	assert.Empty(t, users)
}
