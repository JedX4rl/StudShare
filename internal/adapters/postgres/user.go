package postgres

import (
	"StudShare/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type UserRepository struct {
	db *sql.DB
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT id, email, password_hash, name, surname, phone, is_admin, created_at, rating FROM users WHERE email = $1`
	row := tx.QueryRowContext(ctx, query, email)
	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("failed finding user by email: %w", err)
	}

	var user domain.User
	if err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Surname,
		&user.Phone,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.Rating,
	); err != nil {
		return nil, fmt.Errorf("failed finding user by email: %w", err)
	}

	return &user, tx.Commit()
}

func (u UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT id, email, password_hash, name, surname, phone, is_admin, created_at, rating FROM users WHERE id = $1`
	row := tx.QueryRowContext(ctx, query, id)
	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("failed finding user by id: %w", err)
	}

	var user domain.User
	if err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Surname,
		&user.Phone,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.Rating,
	); err != nil {
		return nil, fmt.Errorf("failed finding user by id: %w", err)
	}
	return &user, tx.Commit()
}

func (u UserRepository) Create(ctx context.Context, user *domain.User) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO users (id, email, password_hash, name, surname, phone)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	if row := tx.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Surname,
		user.Phone,
	); row.Err() != nil {
		return fmt.Errorf("failed creating user: %w", row.Err())
	}
	return tx.Commit()
}

func (u UserRepository) Update(ctx context.Context, updated *domain.User) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
		return fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE users SET email = $1, password_hash = $2, name = $3, surname = $4, phone =  $5 WHERE id = $6`

	if _, err = tx.ExecContext(ctx, query, updated.Email, updated.Password, updated.Name, updated.Surname, updated.Phone, updated.ID); err != nil {
		slog.Error(err.Error())
		return fmt.Errorf("failed updating user: %w", err)
	}
	return tx.Commit()
}

func (u UserRepository) Delete(ctx context.Context, id string) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM users WHERE id = $1`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed deleting user: %w", err)
	}

	return tx.Commit()
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
