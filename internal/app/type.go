package app

import "time"

type UserSubscription struct {
	Id        int
	Service   string
	UserId    string
	Price     uint
	StartDate time.Time
	EndDate   *time.Time
}
