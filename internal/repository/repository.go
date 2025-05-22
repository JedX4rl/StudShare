package repository

import (
	mymongo "StudShare/internal/adapters/mongo"
	"StudShare/internal/adapters/postgres"
	myredis "StudShare/internal/adapters/redis"
	"StudShare/internal/adapters/s3"
	sc "StudShare/internal/config/storage_config"
	"StudShare/internal/domain"
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type ListingRepo interface {
	Create(ctx context.Context, l *domain.Listing) error
	GetByID(ctx context.Context, id string) (*domain.Listing, error)
	FindAll(ctx context.Context, filter string) ([]*domain.Listing, error)
	Update(ctx context.Context, l *domain.Listing) error
	Delete(ctx context.Context, id string) error
	FindNearLocation(ctx context.Context, lat, lon, radius float64, status string) ([]*domain.Listing, error)
}

type ReviewRepo interface {
	Create(ctx context.Context, r *domain.Review) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.Review, error)
	FindByID(ctx context.Context, id string) (*domain.Review, error)
	Delete(ctx context.Context, id string) error
}

type DraftRepo interface {
	Create(ctx context.Context, draft *domain.Draft) error
	FindByID(ctx context.Context, id string) (*domain.Draft, error)
	FindAllByOwner(ctx context.Context, ownerID string) ([]*domain.Draft, error)
	Delete(ctx context.Context, id string) error
}

type FileStorage interface {
	GeneratePresignedPutURL(key string, contentType string, expires time.Duration) (string, error)
	DeleteFile(ctx context.Context, key string) error
	FileExists(ctx context.Context, key string) (bool, error)
	GetBaseURL() string
	GetBucket() string
}

type CacheRepo interface {
	SetUserData(ctx context.Context, user *domain.User) error
	GetUserData(ctx context.Context, userID string) (domain.User, error)
	DeleteUserData(ctx context.Context, userID string) error
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	SaveListing(ctx context.Context, listing *domain.Listing) error
	GetListingByID(ctx context.Context, id string) (*domain.Listing, error)
	DeleteListing(ctx context.Context, id string) error
}

type Repository struct {
	UserRepo
	ListingRepo
	ReviewRepo
	DraftRepo
	FileStorage
	CacheRepo
}

func NewRepository(pg *sql.DB, rd *redis.Client, mg *mongo.Database, s3 *s3.S3, s3Config sc.S3Config) *Repository {
	return &Repository{
		UserRepo:    postgres.NewUserRepo(pg),
		ListingRepo: postgres.NewListingRepo(pg),
		ReviewRepo:  postgres.NewReviewRepo(pg),
		DraftRepo:   mymongo.NewReportRepo(mg),
		CacheRepo:   myredis.NewCacheRepo(rd),
		FileStorage: mys3.NewFileStorage(s3, s3Config),
	}
}
