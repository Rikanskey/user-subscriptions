package query

import (
	"context"
	"github.com/pkg/errors"
	"user-subscriptions/internal/app"
)

type getSubModel interface {
	GetSub(ctx context.Context, subId int) (app.UserSubscription, error)
}

type GetSubHandler struct {
	readModel getSubModel
}

func NewGetSubHandler(readModel getSubModel) GetSubHandler {
	return GetSubHandler{readModel: readModel}
}

func (h GetSubHandler) Handle(ctx context.Context, subId int) (app.UserSubscription, error) {
	sub, err := h.readModel.GetSub(ctx, subId)

	return sub, errors.Wrapf(err, "error getting sub %d", subId)
}
