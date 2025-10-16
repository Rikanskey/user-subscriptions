package query

import (
	"context"
	"github.com/pkg/errors"
	"user-subscriptions/internal/app"
)

type getSubByUserIdModel interface {
	GetSubsByUserId(ctx context.Context, usrId string) ([]app.UserSubscription, error)
}

type GetSubUsrIdHandler struct {
	readModel getSubByUserIdModel
}

func NewGetSubUsrIdHandler(readModel getSubByUserIdModel) GetSubUsrIdHandler {
	return GetSubUsrIdHandler{readModel: readModel}
}

func (h GetSubUsrIdHandler) Handle(ctx context.Context, userId string) ([]app.UserSubscription, error) {
	subs, err := h.readModel.GetSubsByUserId(ctx, userId)

	return subs, errors.Wrapf(err, "error getting sub %d", userId)
}
