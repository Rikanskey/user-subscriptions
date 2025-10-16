package query

import (
	"context"
	"time"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/domain/subs"
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
	) ([]app.UserSubscription, error)
}

type GetSubsUserServiceDate struct {
	readModel getSubsByFilterReadModel
}

func NewGetSubsUserServiceDate(readModel getSubsByFilterReadModel) GetSubsUserServiceDate {
	return GetSubsUserServiceDate{readModel: readModel}
}

func (h GetSubsUserServiceDate) Handle(ctx context.Context, qry app.GetSubByUserServiceDateQuery) (uint, error) {
	subsDb, err := h.readModel.GetSubsByUserIdServiceNameStarDateEndDate(ctx, SubsFilterParams{
		UserId:    qry.UserId,
		Service:   qry.Service,
		StartDate: qry.StartDate,
		EndDate:   qry.EndDate,
	})
	if err != nil {
		return 0, err
	}

	unmarshallingParams := make([]subs.UnmarshallingParams, len(subsDb))

	for _, s := range subsDb {
		unmarshallingParams = append(unmarshallingParams, subs.UnmarshallingParams{
			Id:        s.Id,
			UserId:    s.UserId,
			Service:   s.Service,
			Price:     s.Price,
			StartDate: s.StartDate,
			EndDate:   s.EndDate,
		})
	}

	userSubscriptions := subs.UnmarshalFromDatabaseCollection(unmarshallingParams)

	return userSubscriptions.SumPrice(), nil
}
