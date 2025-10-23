package subs

import (
	"github.com/google/uuid"
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
	Id        int
	Service   string
	Price     uint
	UserId    string
	StartDate time.Time
	EndDate   *time.Time
}

var (
	ErrEmptyUserId   = errors.New("user id cannot be empty")
	ErrEmptyService  = errors.New("service name cannot be empty")
	ErrInvalidUserId = errors.New("invalid user uuid")
)

func isUserIdCorrect(id string) bool {
	return uuid.Validate(id) == nil
}

func IsInvalidSubscriptionParameterError(err error) bool {
	return errors.Is(err, ErrInvalidUserId) || errors.Is(err, ErrEmptyService) || errors.Is(err, ErrEmptyUserId)
}

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
	if !isUserIdCorrect(params.UserId) {
		return UsrSubscription{}, ErrInvalidUserId
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
	if !isUserIdCorrect(params.UserId) {
		return UsrSubscription{}, ErrInvalidUserId
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
