package service

import (
	"StudShare/internal/auth"
	"StudShare/internal/domain"
	"StudShare/internal/repository"
	"context"
)

type Auth interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
}

type Users interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
	GetUserRating(ctx context.Context, userID string) (float64, error)
}

type Listings interface {
	CreateListing(ctx context.Context, listing *domain.Listing) error
	GetListingByID(ctx context.Context, id string) (*domain.Listing, error)
	GetAllListings(ctx context.Context, filter string) ([]*domain.Listing, error)
	UpdateListing(ctx context.Context, listing *domain.Listing) error
	DeleteListing(ctx context.Context, id string) error
	GetListingsNear(ctx context.Context, lat, lon, radius float64, status string) ([]*domain.Listing, error)
}

type Reviews interface {
	AddReview(ctx context.Context, review *domain.Review) error
	GetReviewsForUser(ctx context.Context, userID string) ([]*domain.Review, error)
	DeleteReview(ctx context.Context, id, userID string) error
}

type Drafts interface {
	Create(ctx context.Context, draft *domain.Draft) error
	GetAllDrafts(ctx context.Context, filter string) ([]*domain.Draft, error)
	GetByID(ctx context.Context, id string) (*domain.Draft, error)
	DeleteByID(ctx context.Context, id string) error
}

type Files interface {
	GenerateUploadURL(ctx context.Context, userID string, req domain.FileRequest) (*domain.FileUploadResponse, error)
	DeleteFile(ctx context.Context, userID string, key string) error
}

type Cache interface {
	SetUsername(ctx context.Context, userID string, username string) error
	GetUsername(ctx context.Context, userID string) (string, error)
	DeleteUsername(ctx context.Context, userID string) error
}

type Service struct {
	Auth
	Users
	Listings
	Reviews
	Drafts
	Files
	Cache
}

func NewService(repositories *repository.Repository, tokenManager *auth.TokenManager) *Service {
	return &Service{
		Auth:     NewAuthService(repositories.UserRepo, repositories.CacheRepo, tokenManager),
		Users:    NewUserService(repositories.UserRepo),
		Listings: NewListingService(repositories.ListingRepo, repositories.UserRepo, repositories.CacheRepo),
		Reviews:  NewReviewService(repositories.ReviewRepo, repositories.UserRepo),
		Drafts:   NewDraftService(repositories.DraftRepo),
		Files:    NewFileService(repositories.FileStorage),
		Cache:    NewCacheService(repositories.CacheRepo),
	}
}
