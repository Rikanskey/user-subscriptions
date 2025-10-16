package repository

import "user-subscriptions/internal/app"

func unmarshalUserSubscription(sub usrSubscription) app.UserSubscription {
	return app.UserSubscription{
		Id:        sub.Id,
		Service:   sub.Service,
		UserId:    sub.UsrId,
		Price:     sub.Price,
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
	}
}
