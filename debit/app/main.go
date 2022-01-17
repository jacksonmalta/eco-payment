package app

import (
	"context"
	"fmt"
)

const (
	Withdraw                = "Withdraw"
	InstallmentBuying       = "InstallmentBuying"
	Buying                  = "Buying"
	UnauthorizedTransaction = "unauthorized-transaction"
	UnauthorizedSettlement  = "unauthorized-settlement"
	AuthorizerNotFound      = "authorizer-not-found"
	OperationTypeInvalid    = "operation-type-invalid"
	SettlementFailed        = "settlement-failed"
)

type debit struct {
	log        Logger
	authorizer Authorizer
	settlement Settlement
}

func (a *debit) TransactionWithContext(ctx context.Context, input *TransactionInput) (*TransactionOutput, error) {
	transactionOutput := &TransactionOutput{
		Error: false,
	}

	t := input.OperationType
	if Withdraw != t && InstallmentBuying != t && Buying != t {
		return &TransactionOutput{
			Error:  true,
			Code:   OperationTypeInvalid,
			Detail: "operation type invalid",
		}, nil
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
		OperationType: input.OperationType,
		Amount:        input.Amount * -1,
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

func New(authorizer Authorizer, settlement Settlement, log Logger) Debit {
	return &debit{
		log:        log,
		authorizer: authorizer,
		settlement: settlement,
	}
}
