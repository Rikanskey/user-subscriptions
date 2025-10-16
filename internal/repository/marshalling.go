package repository

import (
	"time"
	"user-subscriptions/internal/domain/subs"
)

type usrSubscription struct {
	Id        int        `db:"id"`
	Service   string     `db:"service"`
	UsrId     string     `db:"usr_id"`
	Price     uint       `db:"price"`
	StartDate time.Time  `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
}

func marshallUsrSubscription(sub subs.UsrSubscription) usrSubscription {
	return usrSubscription{
		Id:        sub.ID(),
		Service:   sub.Service(),
		UsrId:     sub.UserId(),
		Price:     sub.Price(),
		StartDate: sub.StartDate(),
		EndDate:   sub.EndDate(),
	}
}
