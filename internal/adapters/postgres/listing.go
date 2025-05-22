package postgres

import (
	"StudShare/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type ListingRepository struct {
	db *sql.DB
}

func (l *ListingRepository) Create(ctx context.Context, listing *domain.Listing) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertListingQuery := `
		INSERT INTO listings (id, title, description, owner_id, preview_url, location, city, street, status)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8, $9, $10)
	`

	_, err = tx.ExecContext(ctx, insertListingQuery,
		listing.ID,
		listing.Title,
		listing.Description,
		listing.Owner.ID,
		listing.PreviewURL,
		listing.Longitude,
		listing.Latitude,
		listing.City,
		listing.Street,
		listing.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert listing: %w", err)
	}

	insertImageQuery := `
		INSERT INTO listing_images (id, listing_id, url)
		VALUES ($1, $2, $3)
	`

	for _, url := range listing.Images {
		_, err = tx.ExecContext(ctx, insertImageQuery, uuid.New().String(), listing.ID, url)
		if err != nil {
			return fmt.Errorf("failed to insert image: %w", err)
		}
	}

	return tx.Commit()
}

func (l *ListingRepository) GetByID(ctx context.Context, id string) (*domain.Listing, error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	SELECT 
		l.id, l.title, l.description, l.city, l.street, l.status, l.created_at,
		ST_Y(l.location::geometry),
		ST_X(l.location::geometry),
		u.id, u.name, u.surname, u.phone, u.rating
	FROM listings l
	JOIN users u ON l.owner_id = u.id
	WHERE l.id = $1;
	`
	var listing domain.Listing

	if err = tx.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.City,
		&listing.Street,
		&listing.Status,
		&listing.CreatedAt,
		&listing.Latitude,
		&listing.Longitude,
		&listing.Owner.ID,
		&listing.Owner.Name,
		&listing.Owner.Surname,
		&listing.Owner.Phone,
		&listing.Owner.Rating,
	); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	imgQuery := `SELECT url FROM listing_images WHERE listing_id = $1`
	rows, err := tx.QueryContext(ctx, imgQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err = rows.Scan(&url); err != nil {
			return nil, err
		}
		listing.Images = append(listing.Images, url)
	}

	return &listing, tx.Commit()
}

func (l *ListingRepository) FindAll(ctx context.Context, status string) ([]*domain.Listing, error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		SELECT 
			l.id, l.title, l.description, l.preview_url, l.status, 
			l.city, l.created_at, u.id,
			u.name, u.surname, u.rating
		FROM listings l
		JOIN users u ON l.owner_id = u.id
	`

	var args []interface{}
	if status != "" {
		query += "WHERE l.status = $1 ORDER BY l.created_at DESC"
		args = append(args, status)
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying listings: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		var listing domain.Listing
		var owner domain.Owner

		err = rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.PreviewURL,
			&listing.Status,
			&listing.City,
			&listing.CreatedAt,
			&owner.ID,
			&owner.Surname,
			&owner.Name,
			&owner.Rating,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning listing row: %w", err)
		}

		listing.Owner = owner
		listings = append(listings, &listing)
	}

	return listings, tx.Commit()
}

func (l *ListingRepository) Update(ctx context.Context, listing *domain.Listing) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE listings SET 
			title = $1, 
			description = $2, 
			city = $3, 
			street = $4,
			status = $5,
			preview_url = $6,
			location = ST_SetSRID(ST_MakePoint($7, $8), 4326)
		WHERE id = $9;
	`
	if _, err = tx.ExecContext(ctx, query,
		listing.Title,
		listing.Description,
		listing.City,
		listing.Street,
		listing.Status,
		listing.PreviewURL,
		listing.Longitude,
		listing.Latitude,
		listing.ID,
	); err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	if len(listing.Images) > 0 {
		delQuery := `DELETE FROM listing_images WHERE listing_id = $1;`
		if _, err = tx.ExecContext(ctx, delQuery, listing.ID); err != nil {
			return fmt.Errorf("failed to delete listing images: %w", err)
		}

		insertQuery := `INSERT INTO listing_images (id, listing_id, url) VALUES ($1, $2, $3);`
		for _, img := range listing.Images {
			if _, err = tx.ExecContext(ctx, insertQuery, uuid.New().String(), listing.ID, img); err != nil {
				return fmt.Errorf("failed to insert listing image: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (l *ListingRepository) Delete(ctx context.Context, id string) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deleteImagesQuery := `DELETE FROM listing_images WHERE listing_id = $1`
	if _, err = tx.ExecContext(ctx, deleteImagesQuery, id); err != nil {
		return fmt.Errorf("failed to delete listing images: %w", err)
	}

	deleteListingQuery := `DELETE FROM listings WHERE id = $1`
	if _, err = tx.ExecContext(ctx, deleteListingQuery, id); err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	return tx.Commit()
}

func (l *ListingRepository) FindNearLocation(ctx context.Context, lat, lon, radius float64, status string) ([]*domain.Listing, error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		SELECT 
			l.id, l.title, l.description, l.preview_url, l.status, l.created_at, l.city, l.street,
			ST_Distance(l.location, ST_MakePoint($1, $2)::geography) as distance_m,
			u.id, u.name, u.surname, u.rating
		FROM listings l
		JOIN users u ON l.owner_id = u.id
		WHERE ST_DWithin(l.location, ST_MakePoint($1, $2)::geography, $3)
	`

	var args []interface{}
	args = append(args, lon, lat, radius)

	if status != "" {
		query += " AND l.status = $4"
		args = append(args, status)
	}

	query += " ORDER BY distance_m ASC"

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying nearby listings: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		var listing domain.Listing
		var owner domain.Owner

		err = rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.PreviewURL,
			&listing.Status,
			&listing.CreatedAt,
			&listing.City,
			&listing.Street,
			&listing.DistanceM,
			&owner.ID,
			&owner.Name,
			&owner.Surname,
			&owner.Rating,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning nearby listing row: %w", err)
		}

		listing.Owner = owner
		listings = append(listings, &listing)
	}

	return listings, tx.Commit()
}

func NewListingRepo(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}
