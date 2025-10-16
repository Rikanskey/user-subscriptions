package subs

import (
	"github.com/pkg/errors"
	"time"
)

type UsrSubscription struct {
	id        int
	service   string
	price     uint
	userId    string
	startDate time.Time
	endDate   *time.Time
}

type CreationParams struct {
	Id        int
	Service   string
	Price     uint
	UserId    string
	StartDate time.Time
	EndDate   *time.Time
}

type UpdateParams struct {
	Service   string
	Price     uint
	UserId    string
	StartDate time.Time
	EndDate   *time.Time
}

var (
	ErrEmptyUserId  = errors.New("user id cannot be empty")
	ErrEmptyService = errors.New("service name cannot be empty")
)

type UsrSubscriptions []UsrSubscription

func (us *UsrSubscription) ID() int { return us.id }

func (us *UsrSubscription) Service() string { return us.service }

func (us *UsrSubscription) Price() uint { return us.price }

func (us *UsrSubscription) UserId() string { return us.userId }

func (us *UsrSubscription) StartDate() time.Time { return us.startDate }

func (us *UsrSubscription) EndDate() *time.Time { return us.endDate }

func NewUsrSubscription(params CreationParams) (UsrSubscription, error) {
	if params.UserId == "" {
		return UsrSubscription{}, ErrEmptyUserId
	}
	if params.Service == "" {
		return UsrSubscription{}, ErrEmptyService
	}

	return UsrSubscription{
		id:        params.Id,
		service:   params.Service,
		price:     params.Price,
		userId:    params.UserId,
		startDate: params.StartDate,
		endDate:   params.EndDate,
	}, nil
}

func UpdateUsrSubscription(params UpdateParams) (UsrSubscription, error) {
	if params.UserId == "" {
		return UsrSubscription{}, ErrEmptyUserId
	}
	if params.Service == "" {
		return UsrSubscription{}, ErrEmptyService
	}

	return UsrSubscription{
		service:   params.Service,
		price:     params.Price,
		userId:    params.UserId,
		startDate: params.StartDate,
		endDate:   params.EndDate,
	}, nil
}

type UnmarshallingParams struct {
	Id        int
	Service   string
	Price     uint
	UserId    string
	StartDate time.Time
	EndDate   *time.Time
}

func UnmarshalFromDatabaseCollection(params []UnmarshallingParams) UsrSubscriptions {
	usrSubscriptions := make(UsrSubscriptions, len(params))
	for i, param := range params {
		usrSubscriptions[i] = UnmarshalFromDatabase(param)
	}

	return usrSubscriptions
}

func UnmarshalFromDatabase(params UnmarshallingParams) UsrSubscription {
	return UsrSubscription{
		id:        params.Id,
		service:   params.Service,
		price:     params.Price,
		userId:    params.UserId,
		startDate: params.StartDate,
		endDate:   params.EndDate,
	}
}

func (ubs UsrSubscriptions) SumPrice() uint {
	var sum uint
	for _, ub := range ubs {
		sum += ub.price
	}
	return sum
}
