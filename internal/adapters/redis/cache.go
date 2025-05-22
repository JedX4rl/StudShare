package myredis

import (
	"StudShare/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository struct {
	client *redis.Client
}

func (c *CacheRepository) SetUserData(ctx context.Context, user *domain.User) error {
	fields := map[string]any{
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
	}

	fields["phone"] = user.Phone

	key := "user:" + user.ID
	return c.client.HSet(ctx, key, fields, 30*time.Minute).Err()
}

func (c *CacheRepository) GetUserData(ctx context.Context, userID string) (domain.User, error) {
	key := "user:" + userID
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	if len(data) == 0 {
		return domain.User{}, redis.Nil
	}

	return domain.User{
		Email:   data["email"],
		Name:    data["name"],
		Surname: data["surname"],
		Phone:   data["phone"],
	}, nil
}

func (c *CacheRepository) DeleteUserData(ctx context.Context, userID string) error {
	return c.client.Del(ctx, "user:"+userID).Err()
}

func (c *CacheRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return c.client.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func (c *CacheRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := c.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (c *CacheRepository) SaveListing(ctx context.Context, listing *domain.Listing) error {
	data, err := json.Marshal(listing)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "listing:"+listing.ID, data, 10*time.Minute).Err()
}

func (c *CacheRepository) GetListingByID(ctx context.Context, id string) (*domain.Listing, error) {
	data, err := c.client.Get(ctx, "listing:"+id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var listing domain.Listing
	if err = json.Unmarshal([]byte(data), &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}

func (c *CacheRepository) DeleteListing(ctx context.Context, id string) error {
	return c.client.Del(ctx, "listing:"+id).Err()
}

func NewCacheRepo(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}
