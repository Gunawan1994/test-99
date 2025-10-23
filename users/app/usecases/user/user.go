package invoice

import (
	"net/http"
	"strconv"
	"time"
	"user/app/dto"
	"user/app/entity"
	"user/app/repositories"
	listing "user/app/repositories/user"
	"user/app/usecases"
	"user/pkg/response"

	pg "user/pkg/pagination"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type event struct {
	RepoUser repositories.User
	Db       *mongo.Client
}

func New(logger *log.Entry, db *mongo.Client) usecases.User {
	return &event{
		RepoUser: listing.New(logger, db),
	}
}

func (v *event) AddUser(ctx echo.Context, req dto.User) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": req}).Info("repositories: AddUser")

	now := time.Now()
	params := make(map[string]interface{})
	params["name"] = req.Name
	params["created_at"] = now
	params["updated_at"] = now

	var id int64
	id, e = v.RepoUser.CreateUser(ctx, params)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error create user")
		e = response.SetResponse(ctx, http.StatusBadRequest,
			e.Error(), nil, nil, false)
		return
	}

	e = response.SetResponse(ctx, http.StatusOK, "Success", nil, map[string]interface{}{
		"id":         id,
		"name":       req.Name,
		"created_at": now.UnixNano() / 1000,
		"updated_at": now.UnixNano() / 1000,
	}, true)

	return
}

func (v *event) FindAllUser(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": nil}).Info("repositories: FindAllUser")

	meta, e := pg.Pagination(ctx)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error Pagination")
		e = response.SetResponse(ctx, http.StatusBadRequest, "Bad Request", nil, nil, false)
		return
	}

	var total int64
	var data []entity.User
	data, total, e = v.RepoUser.GetAllUser(ctx, meta)
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

func (v *event) FindOneUser(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.WithFields(logrus.Fields{"params": nil}).Info("repositories: FindOneListing")

	id := ctx.Param("id")

	var idInt int64
	idInt, e = strconv.ParseInt(id, 10, 64)
	if e != nil {
		logger.WithField("error", e.Error()).Error("Catch error data")
		e = response.SetResponse(ctx, http.StatusBadRequest,
			e.Error(), nil, nil, false)
		return
	}

	var data entity.User
	data, e = v.RepoUser.GetUser(ctx, idInt)
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
