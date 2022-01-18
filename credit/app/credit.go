package app

import "context"

type Credit interface {
	TransactionWithContext(ctx context.Context, input *TransactionInput) (*TransactionOutput, error)
}

type TransactionInput struct {
	AccountKey  string
	ExternalKey string
	Amount      int
}

type TransactionOutput struct {
	Error  bool
	Code   string
	Detail string
}
