package command

import (
	"context"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/domain/subs"
)

type AddSubHandler struct {
	subsRepository subsRepository
}

func NewAddSubHandler(subsRepository subsRepository) AddSubHandler {
	return AddSubHandler{subsRepository: subsRepository}
}

func (h AddSubHandler) Handle(ctx context.Context, cmd app.AddSubUserCommand) (int, error) {
	sub, err := subs.NewUsrSubscription(subs.CreationParams{
		Service:   cmd.Service,
		UserId:    cmd.UserId,
		Price:     cmd.Price,
		StartDate: cmd.StartDate,
		EndDate:   cmd.EndDate,
	})

	if err != nil {
		return 0, err
	}

	id, err := h.subsRepository.AddSub(ctx, sub)

	return id, err
}
