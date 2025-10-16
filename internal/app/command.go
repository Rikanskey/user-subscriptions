package app

import "time"

type (
	AddSubUserCommand struct {
		Service   string
		UserId    string
		Price     uint
		StartDate time.Time
		EndDate   *time.Time
	}

	UpdateSubCommand struct {
		Id        int
		Service   string
		UserId    string
		Price     uint
		StartDate time.Time
		EndDate   *time.Time
	}
)
