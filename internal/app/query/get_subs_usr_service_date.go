package query

import (
	"context"
	"time"
	"user-subscriptions/internal/app"
)

type SubsFilterParams struct {
	UserId    string
	Service   string
	StartDate time.Time
	EndDate   time.Time
}

type getSubsByFilterReadModel interface {
	GetSubsByUserIdServiceNameStarDateEndDate(
		ctx context.Context,
		sfp SubsFilterParams,
	) (uint, error)
}

type GetSubsUserServiceDate struct {
	readModel getSubsByFilterReadModel
}

func NewGetSubsUserServiceDate(readModel getSubsByFilterReadModel) GetSubsUserServiceDate {
	return GetSubsUserServiceDate{readModel: readModel}
}

func (h GetSubsUserServiceDate) Handle(ctx context.Context, qry app.GetSubByUserServiceDateQuery) (uint, error) {
	return h.readModel.GetSubsByUserIdServiceNameStarDateEndDate(ctx, SubsFilterParams{
		UserId:    qry.UserId,
		Service:   qry.Service,
		StartDate: qry.StartDate,
		EndDate:   qry.EndDate,
	})
}
