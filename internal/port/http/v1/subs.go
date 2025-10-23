package v1

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/domain/subs"
	"user-subscriptions/pkg/httperr"
)

func (h handler) CreateUserSub(w http.ResponseWriter, r *http.Request) {
	cmd, ok := unmarshalUserSubCreateRequest(w, r)
	if !ok {
		return
	}

	id, err := h.app.Commands.AddSubUserCommand.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Content-Location", fmt.Sprintf("/subs/%d", id))
		w.WriteHeader(http.StatusCreated)
		return
	} else if subs.IsInvalidSubscriptionParameterError(err) {
		httperr.UnprocessableEntity("invalid-subscription-parameters", err, w, r)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetSubsFindByUser(w http.ResponseWriter, r *http.Request, params GetSubsFindByUserParams) {
	qry, ok := unmarshallFindByUserRequest(w, r, params)
	if !ok {
		return
	}

	subsRes, err := h.app.Queries.GetSubsByUserId.Handle(r.Context(), qry)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		marshalUserSubscriptionResponses(w, r, subsRes)
		return
	}

	if errors.Is(err, app.ErrUserDoesNotExistOrOutOfPage) {
		httperr.NotFound("user-not-found", err, w, r)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) DeleteSub(w http.ResponseWriter, r *http.Request, subId int) {
	err := h.app.Commands.RemoveSubCommand.Handle(r.Context(), subId)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if errors.Is(err, app.ErrUserDoesNotExist) {
		httperr.NotFound("user-not-found", err, w, r)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetSub(w http.ResponseWriter, r *http.Request, subId int) {
	sub, err := h.app.Queries.GetSub.Handle(r.Context(), subId)
	if err == nil {
		marshalUserSubscription(w, r, sub)
		return
	} else if errors.Is(err, app.ErrUserSubscriptionDoesNotExist) {
		httperr.NotFound("sub-not-found", err, w, r)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) UpdateSub(w http.ResponseWriter, r *http.Request, subId int) {
	cmd, ok := unmarshalUserSubUpdateRequest(w, r, subId)
	if !ok {
		return
	}

	err := h.app.Commands.UpdateSubCommand.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	} else if subs.IsInvalidSubscriptionParameterError(err) {
		httperr.UnprocessableEntity("invalid-subscription-parameters", err, w, r)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetSubsGetSumPrice(w http.ResponseWriter, r *http.Request, params GetSubsGetSumPriceParams) {
	qry, err := unmarshalFindByUserServicePeriodRequest(w, r, params)
	if err != nil {
		return
	}
	price, err := h.app.Queries.GetSubsPrice.Handle(r.Context(), qry)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		marshalUserSubscriptionPriceResponse(w, r, price)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}
