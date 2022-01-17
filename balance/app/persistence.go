package app

import "context"

type Persistence interface {
	InsertWithContext(ctx context.Context, input *InsertInput) (*InsertOutput, error)
}

type InsertInput struct {
	AccountKey     string
	ExternalKey    string
	OperatiionType string
	Amount         int
}
type InsertOutput struct {
	AlreadyExists bool
}
