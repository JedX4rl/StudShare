package service

import (
	"StudShare/internal/domain"
	"StudShare/internal/my_errors"
	"StudShare/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type ReviewService struct {
	reviewRepo repository.ReviewRepo
	userRepo   repository.UserRepo
}

func (r ReviewService) AddReview(c context.Context, review *domain.Review) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	_, err := r.userRepo.FindByID(ctx, review.TargetID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return r.reviewRepo.Create(ctx, review)
}

func (r ReviewService) GetReviewsForUser(c context.Context, userID string) ([]*domain.Review, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	return r.reviewRepo.FindByUserID(ctx, userID)
}

func (r ReviewService) DeleteReview(c context.Context, id, userID string) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	review, err := r.reviewRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_errors.ErrNotFound
		}
		slog.Error(err.Error())
		return err
	}

	if review.AuthorID != userID {
		return my_errors.ErrPermissionDenied
	}

	return r.reviewRepo.Delete(ctx, id)
}

func NewReviewService(reviewRepo repository.ReviewRepo, userRepo repository.UserRepo) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
	}
}
