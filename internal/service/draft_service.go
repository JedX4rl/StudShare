package service

import (
	"StudShare/internal/domain"
	"StudShare/internal/repository"
	"context"
	"github.com/google/uuid"
	"time"
)

type DraftService struct {
	repo repository.DraftRepo
}

func (r *DraftService) Create(c context.Context, draft *domain.Draft) error {
	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()
	draft.ID = uuid.New().String()
	return r.repo.Create(ctx, draft)
}

func (r *DraftService) GetAllDrafts(c context.Context, ownerID string) ([]*domain.Draft, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()
	return r.repo.FindAllByOwner(ctx, ownerID)

}

func (r *DraftService) GetByID(c context.Context, id string) (*domain.Draft, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()
	return r.repo.FindByID(ctx, id)
}

func (r *DraftService) DeleteByID(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()
	return r.repo.Delete(ctx, id)
}

func NewDraftService(reportRepo repository.DraftRepo) *DraftService {
	return &DraftService{repo: reportRepo}
}
