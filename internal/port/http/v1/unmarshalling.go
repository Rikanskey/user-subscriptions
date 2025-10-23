package v1

import (
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"time"
	"user-subscriptions/internal/app"
	"user-subscriptions/pkg/httperr"
)

const dateLayout = "01-2006"

var (
	errPriceLessThanZero = errors.New("price must be greater than zero")
)

func decode(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := render.Decode(r, v); err != nil {
		{
			httperr.BadRequest("bad-request", err, w, r)
			return false
		}
	}

	return true
}

func unmarshalUserSubCreateRequest(w http.ResponseWriter, r *http.Request) (cmd app.AddSubUserCommand, ok bool) {
	var createUserSubRequest CreateUserSubRequest
	if ok = decode(w, r, &createUserSubRequest); !ok {
		return app.AddSubUserCommand{}, false
	}

	if createUserSubRequest.Price < 0 {
		httperr.BadRequest("bad-request", errPriceLessThanZero, w, r)
		return app.AddSubUserCommand{}, false
	}

	startDate, err := time.Parse(dateLayout, createUserSubRequest.StartDate)
	if err != nil {
		httperr.BadRequest("bad-request", err, w, r)
		return app.AddSubUserCommand{}, false
	}

	var ed time.Time
	if createUserSubRequest.EndDate != nil {
		ed, err = time.Parse(dateLayout, *createUserSubRequest.EndDate)
		if err != nil {
			httperr.BadRequest("bad-request", err, w, r)
			return app.AddSubUserCommand{}, false
		}
	} else {
		return app.AddSubUserCommand{
			UserId:    createUserSubRequest.UserId,
			Service:   createUserSubRequest.ServiceName,
			Price:     uint(createUserSubRequest.Price),
			StartDate: startDate,
		}, true
	}

	return app.AddSubUserCommand{
		UserId:    createUserSubRequest.UserId,
		Service:   createUserSubRequest.ServiceName,
		Price:     uint(createUserSubRequest.Price),
		StartDate: startDate,
		EndDate:   &ed,
	}, true
}

func unmarshalUserSubUpdateRequest(w http.ResponseWriter, r *http.Request, subId int) (cmd app.UpdateSubCommand, ok bool) {
	var updateRequest UpdateUserSubRequest
	if ok = decode(w, r, &updateRequest); !ok {
		return
	}

	if updateRequest.Price < 0 {
		httperr.BadRequest("bad-request", errPriceLessThanZero, w, r)
		return app.UpdateSubCommand{}, false
	}

	startDate, err := time.Parse(dateLayout, updateRequest.StartDate)
	if err != nil {
		httperr.BadRequest("bad-request", err, w, r)
		return app.UpdateSubCommand{}, false
	}

	var endDate time.Time
	if updateRequest.EndDate != nil {
		endDate, err = time.Parse(dateLayout, *updateRequest.EndDate)
		if err != nil {
			httperr.BadRequest("bad-request", err, w, r)
			return app.UpdateSubCommand{}, false
		}
	}

	return app.UpdateSubCommand{
		Id:        subId,
		UserId:    updateRequest.UserId,
		Service:   updateRequest.ServiceName,
		Price:     uint(updateRequest.Price),
		StartDate: startDate,
		EndDate:   &endDate,
	}, true
}

func unmarshalFindByUserServicePeriodRequest(w http.ResponseWriter, r *http.Request, params GetSubsGetSumPriceParams) (app.GetSubByUserServiceDateQuery, error) {
	startDate, err := time.Parse(dateLayout, params.StartDate)
	if err != nil {
		httperr.BadRequest("bad-request", err, w, r)
		return app.GetSubByUserServiceDateQuery{}, err
	}
	endDate, err := time.Parse(dateLayout, params.EndDate)
	if err != nil {
		httperr.BadRequest("bad-request", err, w, r)
		return app.GetSubByUserServiceDateQuery{}, err
	}

	return app.GetSubByUserServiceDateQuery{
		UserId:    params.UserId,
		Service:   params.Service,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func unmarshallFindByUserRequest(w http.ResponseWriter, r *http.Request, params GetSubsFindByUserParams) (app.GetSubsByUser, bool) {
	if params.UserId == "" || uuid.Validate(params.UserId) != nil || params.Page <= 0 || params.Limit <= 0 {
		httperr.BadRequest("incorrect-user-id", app.ErrUserDoesNotExist, w, r)
		return app.GetSubsByUser{}, false
	}

	return app.GetSubsByUser{UserId: params.UserId, Page: params.Page, Limit: params.Limit}, true
}
