package app

import "context"

type Balance interface {
	SettlementWithContext(ctx context.Context, input *SettlementInput) (*SettlementOutput, error)
}

type SettlementInput struct {
	AccountKey    string
	ExternalKey   string
	OperationType string
	Amount        int
}

type SettlementOutput struct {
	Error  bool
	Code   string
	Detail string
}
