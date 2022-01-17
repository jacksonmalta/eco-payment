package app

import (
	"context"
	"fmt"
)

const (
	Payment                 = "Payment"
	UnauthorizedTransaction = "unauthorized-transaction"
	UnauthorizedSettlement  = "unauthorized-settlement"
	AuthorizerNotFound      = "authorizer-not-found"
	SettlementFailed        = "settlement-failed"
)

type credit struct {
	log        Logger
	authorizer Authorizer
	settlement Settlement
}

func (a *credit) TransactionWithContext(ctx context.Context, input *TransactionInput) (*TransactionOutput, error) {
	transactionOutput := &TransactionOutput{
		Error: false,
	}

	ai := &AuthorizeInput{
		AccountKey: input.AccountKey,
	}
	ao, err := a.authorizer.AuthorizeWithContext(ctx, ai)
	if err != nil {
		a.log.Error(fmt.Sprintf("authorize error %s", err.Error()))
		return nil, err
	}
	if ao == nil {
		return &TransactionOutput{
			Error:  true,
			Code:   AuthorizerNotFound,
			Detail: "authorizer not found",
		}, nil
	}
	if ao.HasError {
		return &TransactionOutput{
			Error:  true,
			Code:   UnauthorizedTransaction,
			Detail: "Try again",
		}, nil
	}

	si := &SettleInput{
		AccountKey:    input.AccountKey,
		ExternalKey:   input.ExternalKey,
		OperationType: Payment,
		Amount:        input.Amount,
	}
	so, err := a.settlement.SettleWithContext(ctx, si)
	if err != nil {
		a.log.Error(fmt.Sprintf("settle error %s", err.Error()))
		return nil, err
	}
	if so.HasIntermitance {
		return &TransactionOutput{
			Error:  true,
			Code:   UnauthorizedSettlement,
			Detail: "Try again",
		}, nil
	}
	if so.Error {
		return &TransactionOutput{
			Error:  true,
			Code:   SettlementFailed,
			Detail: so.Detail,
		}, nil
	}

	return transactionOutput, nil
}

func New(authorizer Authorizer, settlement Settlement, log Logger) Credit {
	return &credit{
		log:        log,
		authorizer: authorizer,
		settlement: settlement,
	}
}
