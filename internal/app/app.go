package app

import "context"

type Application struct {
	Commands Commands
	Queries  Queries
}

type (
	Commands struct {
		AddSubUserCommand addSubUserCommand
		RemoveSubCommand  removeSubCommand
		UpdateSubCommand  updateSubCommand
	}
	addSubUserCommand interface {
		Handle(ctx context.Context, addSubUserCommand AddSubUserCommand) (int, error)
	}
	removeSubCommand interface {
		Handle(ctx context.Context, subId int) error
	}
	updateSubCommand interface {
		Handle(ctx context.Context, cmd UpdateSubCommand) error
	}
)

type (
	Queries struct {
		GetSub          getSub
		GetSubsByUserId getSubByUserId
		GetSubsPrice    getSubsPriceByUserIdServiceDate
	}
	getSub interface {
		Handle(ctx context.Context, subId int) (UserSubscription, error)
	}
	getSubByUserId interface {
		Handle(ctx context.Context, userId string) ([]UserSubscription, error)
	}
	getSubsPriceByUserIdServiceDate interface {
		Handle(ctx context.Context, qry GetSubByUserServiceDateQuery) (uint, error)
	}
)
