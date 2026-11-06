package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"go-zhaboilerplate/domain"
	"go-zhaboilerplate/internal/repository"
)

type UserRepository struct {
	Conn *sql.DB
}

// NewUserRepository creates an instance of UserRepository
func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{Conn: conn}
}

func (m *UserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.User, 0)
	for rows.Next() {
		u := domain.User{}
		err = rows.Scan(
			&u.ID,
			&u.Email,
			&u.Password,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, u)
	}

	return result, nil
}

func (m *UserRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	query := `SELECT id, email, password, is_active, created_at, updated_at
				FROM users WHERE created_at > ? ORDER BY created_at LIMIT ?`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, nil
}

func (m *UserRepository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
	query := `SELECT id, email, password, is_active, created_at, updated_at FROM users WHERE id = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return domain.User{}, domain.ErrNotFound
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	query := `SELECT id, email, password, is_active, created_at, updated_at FROM users WHERE email = ?`

	list, err := m.fetch(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return domain.User{}, domain.ErrNotFound
}

func (m *UserRepository) Store(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (email, password, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Email, u.Password, u.IsActive, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = lastID
	return nil
}

func (m *UserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET email = ?, password = ?, is_active = ?, updated_at = ? WHERE id = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Email, u.Password, u.IsActive, u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (m *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (m *UserRepository) SetActive(ctx context.Context, id int64, isActive bool) error {
	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, isActive, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

