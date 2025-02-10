package repositories

import (
	"context"
	// "database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"proposal-template/models"
)

// TestNewUserRepo verifies that UserRepo is initialized correctly
func TestNewUserRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.GenericDAO)
}

// TestGetByColumn tests retrieving a user by a specific column (e.g., email)
func TestGetByColumn(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	ctx := context.Background()

	// Create a test UUID and timestamp
	testUUID := uuid.New()
	testCreatedAt := time.Now()
	testUpdatedAt := time.Now()

	// Expected user record
	expectedUser := model.User{
		BaseModel: model.BaseModel{
			UUID:      testUUID,
			CreatedAt: testCreatedAt,
			UpdatedAt: testUpdatedAt,
		},
		Name:          "John Doe",
		Email:         "johndoe@example.com",
		EmailVerified: true,
	}

	// Mock SQL query response
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "email", "email_verified"}).
		AddRow(testUUID, testCreatedAt, testUpdatedAt, expectedUser.Name, expectedUser.Email, expectedUser.EmailVerified)

	mock.ExpectQuery(`SELECT \* FROM users WHERE email = \$1 LIMIT 1`).
		WithArgs(expectedUser.Email).
		WillReturnRows(rows)

	// Call GetByColumn
	user, err := repo.GetByColumn(ctx, "email", expectedUser.Email)
	// assert.NoError(t, err)
	// assert.NotNil(t, user)
	assert.Equal(t, expectedUser.UUID, user.UUID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.EmailVerified, user.EmailVerified)
	assert.WithinDuration(t, expectedUser.CreatedAt, user.CreatedAt, time.Second)
	assert.WithinDuration(t, expectedUser.UpdatedAt, user.UpdatedAt, time.Second)

	// Ensure expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestList tests retrieving multiple users with pagination
func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	ctx := context.Background()
	paging := Paging{Page: 1, Limit: 2}

	// Create test data
	testUUID1 := uuid.New()
	testUUID2 := uuid.New()
	testCreatedAt := time.Now()
	testUpdatedAt := time.Now()

	// Mock SQL query result
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "email", "email_verified"}).
		AddRow(testUUID1, testCreatedAt, testUpdatedAt, "John Doe", "john@example.com", true).
		AddRow(testUUID2, testCreatedAt, testUpdatedAt, "Jane Doe", "jane@example.com", false)

	mock.ExpectQuery(`SELECT \* FROM users ORDER BY id DESC LIMIT 2 OFFSET 0`).
		WillReturnRows(rows)

	// Call List
	users, err := repo.List(ctx, paging, "SELECT * FROM users")
	// assert.NoError(t, err)
	// assert.Len(t, users, 2)

	// Check first user
	assert.Equal(t, "John Doe", users[0].Name)
	assert.Equal(t, "john@example.com", users[0].Email)
	assert.Equal(t, true, users[0].EmailVerified)
	assert.Equal(t, testUUID1, users[0].UUID)

	// Check second user
	assert.Equal(t, "Jane Doe", users[1].Name)
	assert.Equal(t, "jane@example.com", users[1].Email)
	assert.Equal(t, false, users[1].EmailVerified)
	assert.Equal(t, testUUID2, users[1].UUID)

	// Ensure expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCreate tests inserting a new user into the database
func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	ctx := context.Background()

	newUser := model.User{
		Name:          "New User",
		Email:         "newuser@example.com",
		EmailVerified: false,
	}

	testUUID := uuid.New()
	testCreatedAt := time.Now()
	testUpdatedAt := time.Now()

	// Mock INSERT query and return inserted UUID
	mock.ExpectQuery(`INSERT INTO users \(.+\) VALUES \(.+\) RETURNING id, created_at, updated_at`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(testUUID, testCreatedAt, testUpdatedAt))

	// Call Create
	id, err := repo.Create(ctx, newUser)
	assert.NoError(t, err)
	assert.Equal(t, testUUID, id)

	// Ensure expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
