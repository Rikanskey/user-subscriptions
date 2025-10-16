package app

import "time"

type GetSubByUserServiceDateQuery struct {
	UserId    string
	Service   string
	StartDate time.Time
	EndDate   time.Time
}
