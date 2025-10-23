package listing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"listing/app/entity"
	"listing/app/repositories"
	"listing/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type event struct {
	Logger *log.Entry
	Db     *sql.DB
}

func New(logger *log.Entry, db *sql.DB) repositories.Listing {
	return &event{
		Logger: logger,
		Db:     db,
	}
}

func (v *event) CreateListing(c echo.Context, paramsListing map[string]interface{}) (id int64, e error) {
	logger := c.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"paramsListing": paramsListing}).Info("repositories: CreateListing")

	ctx := context.Background()
	tx, e := v.Db.BeginTx(ctx, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO listings("
	var fields []string
	var placeholders []string
	var values []interface{}

	i := 1
	for key, val := range paramsListing {
		fields = append(fields, fmt.Sprintf("\"%s\"", key))
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, val)
		i++
	}

	query += fmt.Sprintf("%s) VALUES(%s) RETURNING id",
		utils.Join(fields, ", "),
		utils.Join(placeholders, ", "),
	)

	logger.WithField("query", query).Info("Executing query")

	e = tx.QueryRowContext(ctx, query, values...).Scan(&id)
	if e != nil {
		logger.WithError(e).Error("Failed to insert listing")
		return
	}

	// Commit transaction
	if e = tx.Commit(); e != nil {
		logger.WithError(e).Error("Failed to commit transaction")
		return
	}

	logger.WithField("inserted_id", id).Info("Listing created successfully")
	return
}

func (v *event) GetAllListing(c echo.Context, meta entity.Meta) (data []entity.Listing, total int64, e error) {
	logger := c.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": meta}).Info("repositories: GetAllListing")

	var (
		createdAt time.Time
		updatedAt time.Time
	)

	ctx := context.Background()

	query := `
		SELECT id, user_id, listing_type, price, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, e := v.Db.QueryContext(ctx, query, meta.Limit, meta.Offset)
	if e != nil {
		logger.WithError(e).Error("failed to query listings")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d entity.Listing
		if e = rows.Scan(
			&d.ID,
			&d.UserID,
			&d.ListingType,
			&d.Price,
			&createdAt,
			&updatedAt,
		); e != nil {
			logger.WithError(e).Error("failed to scan listing row")
			return
		}

		d.CreatedAt = createdAt.Unix()
		d.UpdatedAt = updatedAt.Unix()
		data = append(data, d)
	}

	countQuery := `SELECT COUNT(*) FROM listings`
	e = v.Db.QueryRowContext(ctx, countQuery).Scan(&total)
	if e != nil {
		logger.WithError(e).Error("failed to count listings")
		return
	}

	return
}

func (v *event) GetListing(c echo.Context, id string) (data entity.Listing, e error) {
	logger := c.Get("logger").(*logrus.Entry)
	logger.WithField("id", id).Info("repositories: GetListing")

	var (
		createdAt time.Time
		updatedAt time.Time
	)

	ctx := context.Background()

	query := `
		SELECT id, user_id, listing_type, price, created_at, updated_at
		FROM listings
		WHERE id = $1
	`

	row := v.Db.QueryRowContext(ctx, query, id)

	e = row.Scan(
		&data.ID,
		&data.UserID,
		&data.ListingType,
		&data.Price,
		&createdAt,
		&updatedAt,
	)
	data.CreatedAt = createdAt.Unix()
	data.UpdatedAt = updatedAt.Unix()
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			logger.Warn("listing not found")
			e = fmt.Errorf("listing not found with id %d", id)
			return
		}
		logger.WithError(e).Error("failed to scan listing")
		return
	}

	return
}
