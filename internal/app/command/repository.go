package command

import (
	"context"
	"user-subscriptions/internal/domain/subs"
)

type subsRepository interface {
	AddSub(ctx context.Context, sub subs.UsrSubscription) (int, error)
	UpdateSub(ctx context.Context, sub subs.UsrSubscription) error
	RemoveSub(ctx context.Context, subId int) error
}
