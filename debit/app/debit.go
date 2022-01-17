package app

import "context"

type Debit interface {
	TransactionWithContext(ctx context.Context, input *TransactionInput) (*TransactionOutput, error)
}

type TransactionInput struct {
	AccountKey    string
	ExternalKey   string
	OperationType string
	Amount        int
}

type TransactionOutput struct {
	Error  bool
	Code   string
	Detail string
}
