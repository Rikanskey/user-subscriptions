package query

import (
	"context"
	"github.com/pkg/errors"
	"user-subscriptions/internal/app"
)

type getSubByUserIdModel interface {
	GetSubsByUserId(ctx context.Context, userId string, page, limit int) ([]app.UserSubscription, error)
}

type GetSubUsrIdHandler struct {
	readModel getSubByUserIdModel
}

func NewGetSubUsrIdHandler(readModel getSubByUserIdModel) GetSubUsrIdHandler {
	return GetSubUsrIdHandler{readModel: readModel}
}

func (h GetSubUsrIdHandler) Handle(ctx context.Context, qry app.GetSubsByUser) ([]app.UserSubscription, error) {
	subs, err := h.readModel.GetSubsByUserId(ctx, qry.UserId, qry.Page, qry.Limit)

	return subs, errors.Wrapf(err, "error getting sub %d", qry.UserId)
}
