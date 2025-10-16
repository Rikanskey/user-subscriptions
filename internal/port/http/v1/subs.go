package v1

import (
	"fmt"
	"net/http"
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
	} else {
		httperr.UnprocessableEntity("invalid-subscription-parameters", err, w, r)
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetSubsFindByUser(w http.ResponseWriter, r *http.Request, params GetSubsFindByUserParams) {
	if params.UserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subs, err := h.app.Queries.GetSubsByUserId.Handle(r.Context(), params.UserId)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		marshalUserSubscriptionResponses(w, r, subs)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) DeleteSub(w http.ResponseWriter, r *http.Request, subId int) {
	err := h.app.Commands.RemoveSubCommand.Handle(r.Context(), subId)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetSub(w http.ResponseWriter, r *http.Request, subId int) {
	sub, err := h.app.Queries.GetSub.Handle(r.Context(), subId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	marshalUserSubscription(w, r, sub)

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
	} else if err != nil {
		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetSubsFindByUserServicePeriod(w http.ResponseWriter, r *http.Request, params GetSubsFindByUserServicePeriodParams) {
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
