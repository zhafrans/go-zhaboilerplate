package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/repository"
	userMysqlRepo "go-zhaboilerplate/internal/repository/mysql"
)

func TestFetchUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	now := time.Now()
	mockUsers := []domain.User{
		{ID: 1, Email: "user1@example.com", Password: "hash1", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Email: "user2@example.com", Password: "hash2", IsActive: false, CreatedAt: now, UpdatedAt: now},
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_active", "created_at", "updated_at"}).
		AddRow(mockUsers[0].ID, mockUsers[0].Email, mockUsers[0].Password, mockUsers[0].IsActive, mockUsers[0].CreatedAt, mockUsers[0].UpdatedAt).
		AddRow(mockUsers[1].ID, mockUsers[1].Email, mockUsers[1].Password, mockUsers[1].IsActive, mockUsers[1].CreatedAt, mockUsers[1].UpdatedAt)

	query := "SELECT id, email, password, is_active, created_at, updated_at FROM users WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	repo := userMysqlRepo.NewUserRepository(db)
	cursor := repository.EncodeCursor(mockUsers[1].CreatedAt)
	list, nextCursor, err := repo.Fetch(context.TODO(), cursor, 2)

	assert.NoError(t, err)
	assert.NotEmpty(t, nextCursor)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	now := time.Now()
	mockUser := domain.User{ID: 1, Email: "user1@example.com", Password: "hash1", IsActive: true, CreatedAt: now, UpdatedAt: now}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_active", "created_at", "updated_at"}).
		AddRow(mockUser.ID, mockUser.Email, mockUser.Password, mockUser.IsActive, mockUser.CreatedAt, mockUser.UpdatedAt)

	query := "SELECT id, email, password, is_active, created_at, updated_at FROM users WHERE id = \\?"

	mock.ExpectQuery(query).WithArgs(int64(1)).WillReturnRows(rows)
	repo := userMysqlRepo.NewUserRepository(db)
	res, err := repo.GetByID(context.TODO(), 1)

	assert.NoError(t, err)
	assert.Equal(t, mockUser.Email, res.Email)
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	now := time.Now()
	mockUser := domain.User{ID: 1, Email: "user1@example.com", Password: "hash1", IsActive: true, CreatedAt: now, UpdatedAt: now}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_active", "created_at", "updated_at"}).
		AddRow(mockUser.ID, mockUser.Email, mockUser.Password, mockUser.IsActive, mockUser.CreatedAt, mockUser.UpdatedAt)

	query := "SELECT id, email, password, is_active, created_at, updated_at FROM users WHERE email = \\?"

	mock.ExpectQuery(query).WithArgs("user1@example.com").WillReturnRows(rows)
	repo := userMysqlRepo.NewUserRepository(db)
	res, err := repo.GetByEmail(context.TODO(), "user1@example.com")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)
}

func TestStoreUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	now := time.Now()
	u := &domain.User{
		Email:     "new@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	query := "INSERT INTO users \\(email, password, is_active, created_at, updated_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.Email, u.Password, u.IsActive, u.CreatedAt, u.UpdatedAt).WillReturnResult(sqlmock.NewResult(10, 1))

	repo := userMysqlRepo.NewUserRepository(db)
	err = repo.Store(context.TODO(), u)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), u.ID)
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	now := time.Now()
	u := &domain.User{
		ID:        5,
		Email:     "update@example.com",
		Password:  "newhash",
		IsActive:  true,
		UpdatedAt: now,
	}

	query := "UPDATE users SET email = \\?, password = \\?, is_active = \\?, updated_at = \\? WHERE id = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.Email, u.Password, u.IsActive, u.UpdatedAt, u.ID).WillReturnResult(sqlmock.NewResult(5, 1))

	repo := userMysqlRepo.NewUserRepository(db)
	err = repo.Update(context.TODO(), u)

	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	query := "DELETE FROM users WHERE id = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(int64(5)).WillReturnResult(sqlmock.NewResult(5, 1))

	repo := userMysqlRepo.NewUserRepository(db)
	err = repo.Delete(context.TODO(), 5)

	assert.NoError(t, err)
}

func TestSetActiveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	query := "UPDATE users SET is_active = \\?, updated_at = \\? WHERE id = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(false, sqlmock.AnyArg(), int64(5)).WillReturnResult(sqlmock.NewResult(5, 1))

	repo := userMysqlRepo.NewUserRepository(db)
	err = repo.SetActive(context.TODO(), 5, false)

	assert.NoError(t, err)
}

