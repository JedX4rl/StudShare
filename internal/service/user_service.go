package service

import (
	"StudShare/internal/domain"
	"StudShare/internal/my_errors"
	"StudShare/internal/repository"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	repo repository.UserRepo
}

func (u *UserService) GetProfile(c context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot find user")
	}
	return &domain.User{
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		Phone:     user.Phone,
		Rating:    user.Rating,
		CreatedAt: user.CreatedAt,
		IsAdmin:   user.IsAdmin,
	}, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, updated *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	currentUser, err := u.repo.FindByID(ctx, updated.ID)
	if err != nil {
		return my_errors.ErrNotFound
	}

	changed := false

	if len(updated.Email) != 0 && currentUser.Email != updated.Email {
		currentUser.Email = updated.Email
		changed = true
	}
	if len(updated.Name) != 0 && currentUser.Name != updated.Name {
		currentUser.Name = updated.Name
		changed = true
	}
	if len(updated.Surname) != 0 && currentUser.Surname != updated.Surname {
		currentUser.Surname = updated.Surname
		changed = true
	}
	if len(updated.Phone) != 0 && currentUser.Phone != updated.Phone {
		currentUser.Phone = updated.Phone
		changed = true
	}
	if updated.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updated.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("bcrypt hashing failed: %w", err)
		}
		currentUser.Password = string(hashedPassword)
		changed = true
	}

	if !changed {
		return my_errors.ErrNoChanges
	}

	return u.repo.Update(ctx, currentUser)
}

func (u *UserService) GetUserRating(ctx context.Context, userID string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{repo: userRepo}
}
