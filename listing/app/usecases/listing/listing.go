package invoice

import (
	"database/sql"
	"listing/app/dto"
	"listing/app/entity"
	"listing/app/repositories"
	"listing/app/repositories/listing"
	"listing/app/usecases"
	"listing/pkg/response"
	"net/http"
	"time"

	pg "listing/pkg/pagination"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type event struct {
	RepoListing repositories.Listing
	AuroraDb    *sql.DB
}

func New(logger *log.Entry, auroraDb *sql.DB) usecases.Listing {
	return &event{
		RepoListing: listing.New(logger, auroraDb),
	}
}

func (v *event) AddListing(ctx echo.Context, req dto.Listing) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": req}).Info("repositories: AddListing")

	now := time.Now()
	params := make(map[string]interface{})
	params["user_id"] = req.UserID
	params["listing_type"] = req.ListingType
	params["price"] = req.Price
	params["created_at"] = now
	params["updated_at"] = now

	var id int64
	id, e = v.RepoListing.CreateListing(ctx, params)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error create listing")
		e = response.SetResponse(ctx, http.StatusBadRequest,
			e.Error(), nil, nil, false)
		return
	}

	e = response.SetResponse(ctx, http.StatusOK, "Success", nil, map[string]interface{}{
		"id":           id,
		"user_id":      req.UserID,
		"listing_type": req.ListingType,
		"price":        req.Price,
		"created_at":   now.UnixNano() / 1000,
		"updated_at":   now.UnixNano() / 1000,
	}, true)

	return
}

func (v *event) FindAllListing(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": nil}).Info("repositories: FindAllListing")

	meta, e := pg.Pagination(ctx)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error Pagination")
		e = response.SetResponse(ctx, http.StatusBadRequest, "Bad Request", nil, nil, false)
		return
	}

	var total int64
	var data []entity.Listing
	data, total, e = v.RepoListing.GetAllListing(ctx, meta)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error data")
		e = response.SetResponse(ctx, http.StatusBadRequest,
			e.Error(), nil, nil, false)
		return
	}
	if len(data) == 0 {
		e = response.SetResponse(ctx, http.StatusNotFound, "Data not found", nil, nil, false)
		return
	}

	metaPagination := pg.GenerateMeta(ctx, total, meta.Limit, meta.Page, meta.Offset, true, nil)

	e = response.SetResponse(ctx, http.StatusOK, "Success", metaPagination, data, true)

	return
}

func (v *event) FindOneListing(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": nil}).Info("repositories: FindOneListing")

	id := ctx.Param("id")

	var data entity.Listing
	data, e = v.RepoListing.GetListing(ctx, id)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error data")
		e = response.SetResponse(ctx, http.StatusBadRequest,
			e.Error(), nil, nil, false)
		return
	}

	if data.ID == 0 {
		e = response.SetResponse(ctx, http.StatusNotFound, "Data not found", nil, nil, false)
		return
	}

	e = response.SetResponse(ctx, http.StatusOK, "Success", nil, data, true)

	return
}
