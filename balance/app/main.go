package app

import (
	"context"
	"fmt"
)

type accreditation struct {
	log        Logger
	repository Persistence
}

func (a *accreditation) SettlementWithContext(ctx context.Context, input *SettlementInput) (*SettlementOutput, error) {
	createAccountOuput := &SettlementOutput{
		Error: false,
	}

	i := &InsertInput{
		AccountKey:     input.AccountKey,
		ExternalKey:    input.ExternalKey,
		OperatiionType: input.OperationType,
		Amount:         input.Amount,
	}

	res, err := a.repository.InsertWithContext(ctx, i)

	if err != nil {
		a.log.Error(fmt.Sprintf("Repository insert error %s", err.Error()))
		return nil, err
	}

	if res != nil && res.AlreadyExists {
		return &SettlementOutput{
			Error:  true,
			Code:   "item-already-exists",
			Detail: "item already exists",
		}, nil
	}

	return createAccountOuput, nil
}

func New(r Persistence, log Logger) Balance {
	return &accreditation{
		repository: r,
		log:        log,
	}
}
