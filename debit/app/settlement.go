package app

import "context"

type Settlement interface {
	SettleWithContext(ctx context.Context, input *SettleInput) (*SettleOutput, error)
}

type SettleInput struct {
	AccountKey    string
	ExternalKey   string
	OperationType string
	Amount        int
}
type SettleOutput struct {
	HasIntermitance bool
	Error           bool
	Code            string
	Detail          string
}
