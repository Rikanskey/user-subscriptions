package command

import "context"

type RemoveSubHandler struct {
	subsRepository subsRepository
}

func NewRemoveSubHandler(subsRepository subsRepository) RemoveSubHandler {
	return RemoveSubHandler{subsRepository: subsRepository}
}

func (h RemoveSubHandler) Handle(ctx context.Context, subId int) error {
	return h.subsRepository.RemoveSub(ctx, subId)
}
