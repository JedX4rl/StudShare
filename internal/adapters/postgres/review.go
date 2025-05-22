package postgres

import (
	"StudShare/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type ReviewRepository struct {
	db *sql.DB
}

func (r *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO reviews (id, author_id, target_id, text, rating)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err = tx.ExecContext(ctx, query, review.ID, review.AuthorID, review.TargetID, review.Comment, review.Rating); err != nil {
		return fmt.Errorf("failed to insert review: %w", err)
	}

	return tx.Commit()
}

func (r *ReviewRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Review, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		SELECT r.id, r.author_id, u.name, r.target_id,
		       r.rating, r.text, r.created_at
		FROM reviews r
		JOIN users u ON u.id = r.author_id
		WHERE r.target_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var review domain.Review

		err = rows.Scan(&review.ID, &review.AuthorID, &review.AuthorName, &review.TargetID,
			&review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return reviews, tx.Commit()
}

func (r *ReviewRepository) FindByID(ctx context.Context, id string) (*domain.Review, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		SELECT r.id, r.author_id, u.name, r.target_id,
		       r.rating, r.text, r.created_at
		FROM reviews r
		JOIN users u ON u.id = r.author_id
		WHERE r.id = $1
		ORDER BY r.created_at DESC
	`

	var review domain.Review

	if err = tx.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.AuthorID,
		&review.AuthorName,
		&review.TargetID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &review, tx.Commit()
}

func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM reviews WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func NewReviewRepo(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}
