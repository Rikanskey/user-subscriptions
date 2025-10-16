package command

import (
	"context"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/domain/subs"
)

type UpdateSubHandler struct {
	subsRepository subsRepository
}

func NewUpdateSubHandler(subsRepository subsRepository) UpdateSubHandler {
	return UpdateSubHandler{subsRepository: subsRepository}
}

func (us UpdateSubHandler) Handle(ctx context.Context, cmd app.UpdateSubCommand) error {
	sub, err := subs.UpdateUsrSubscription(subs.UpdateParams{
		UserId:    cmd.UserId,
		Service:   cmd.Service,
		Price:     cmd.Price,
		StartDate: cmd.StartDate,
		EndDate:   cmd.EndDate,
	})
	if err != nil {
		return err
	}
	return us.subsRepository.UpdateSub(ctx, sub)
}
