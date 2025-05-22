package mymongo

import (
	"StudShare/internal/domain"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DraftRepository struct {
	collection *mongo.Collection
}

func (r DraftRepository) Delete(ctx context.Context, id string) error {

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r DraftRepository) FindByID(ctx context.Context, id string) (*domain.Draft, error) {
	var draft domain.Draft
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&draft)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("draft with id %s not found", id)
		}
		return nil, err
	}
	return &draft, nil
}

func (r DraftRepository) FindAllByOwner(ctx context.Context, ownerID string) ([]*domain.Draft, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drafts []*domain.Draft
	if err = cursor.All(ctx, &drafts); err != nil {
		return nil, err
	}

	return drafts, nil
}

func (r DraftRepository) Create(ctx context.Context, draft *domain.Draft) error {
	_, err := r.collection.InsertOne(ctx, draft)
	return err
}

func NewReportRepo(db *mongo.Database) *DraftRepository {
	return &DraftRepository{
		collection: db.Collection("reports"),
	}
}
