package service

import (
	"StudShare/internal/domain"
	"StudShare/internal/my_errors"
	"StudShare/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"slices"
	"time"
)

type ListingService struct {
	userRepo  repository.UserRepo
	repo      repository.ListingRepo
	cacheRepo repository.CacheRepo
}

func (l *ListingService) CreateListing(c context.Context, listing *domain.Listing) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	owner, err := l.userRepo.FindByID(ctx, listing.Owner.ID)
	if err != nil {
		return fmt.Errorf("error finding owner by id: %w", err)
	}
	listing.Owner.Name = owner.Name
	listing.Owner.Surname = owner.Surname
	listing.Owner.Phone = owner.Phone
	listing.Owner.Rating = owner.Rating
	listing.ID = uuid.New().String()

	return l.repo.Create(ctx, listing)
}

func (l *ListingService) GetListingByID(c context.Context, id string) (*domain.Listing, error) {
	ctx, cancel := context.WithTimeout(c, 50*time.Second)
	defer cancel()

	listing, err := l.cacheRepo.GetListingByID(ctx, id)
	if err == nil && listing != nil {
		return listing, nil
	}

	listing, err = l.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get listing from DB: %w", err)
	}

	_ = l.cacheRepo.SaveListing(ctx, listing)

	return listing, nil
}

func (l *ListingService) GetAllListings(ctx context.Context, filter string) ([]*domain.Listing, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return l.repo.FindAll(ctx, filter)
}

func (l *ListingService) UpdateListing(c context.Context, updated *domain.Listing) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	current, err := l.repo.GetByID(ctx, updated.ID)
	if err != nil {
		return my_errors.ErrNotFound
	}

	changed := false

	if updated.Title != "" && updated.Title != current.Title {
		current.Title = updated.Title
		changed = true
	}
	if updated.Description != "" && updated.Description != current.Description {
		current.Description = updated.Description
		changed = true
	}
	if updated.City != "" && updated.City != current.City {
		current.City = updated.City
		changed = true
	}
	if updated.Street != "" && updated.Street != current.Street {
		current.Street = updated.Street
		changed = true
	}
	if updated.Status != "" && updated.Status != current.Status {
		current.Status = updated.Status
		changed = true
	}
	if updated.PreviewURL != "" && updated.PreviewURL != current.PreviewURL {
		current.PreviewURL = updated.PreviewURL
		changed = true
	}
	if updated.Latitude != 0 && updated.Latitude != current.Latitude {
		current.Latitude = updated.Latitude
		changed = true
	}
	if updated.Longitude != 0 && updated.Longitude != current.Longitude {
		current.Longitude = updated.Longitude
		changed = true
	}
	if len(updated.Images) > 0 && !slices.Equal(updated.Images, current.Images) {
		current.Images = updated.Images
		changed = true
	}

	if !changed {
		return my_errors.ErrNoChanges
	}

	if err = l.repo.Update(ctx, current); err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	_ = l.cacheRepo.DeleteListing(ctx, current.ID)

	return nil
}

func (l *ListingService) DeleteListing(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	userID := ctx.Value("userID").(string)

	listing, err := l.GetListingByID(ctx, id)
	if err != nil {
		return err
	}
	if userID != listing.Owner.ID {
		return my_errors.ErrPermissionDenied
	}
	if err = l.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	_ = l.cacheRepo.DeleteListing(ctx, id)

	return nil
}

func (l *ListingService) GetListingsNear(ctx context.Context, lat, lon, radius float64, status string) ([]*domain.Listing, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return l.repo.FindNearLocation(ctx, lat, lon, radius, status)
}

func NewListingService(listingRepo repository.ListingRepo, userRepo repository.UserRepo, cache repository.CacheRepo) *ListingService {
	return &ListingService{
		repo:      listingRepo,
		userRepo:  userRepo,
		cacheRepo: cache,
	}
}
