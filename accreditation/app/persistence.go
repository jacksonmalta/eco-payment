package app

import "context"

type Persistence interface {
	InsertWithContext(ctx context.Context, input *InsertInput) (*InsertOutput, error)
	GetWithContext(ctx context.Context, input *GetInput) (*GetOutput, error)
}

type InsertInput struct {
	ExternalKey    string
	DocumentNumber string
}
type InsertOutput struct {
	AlreadyExists bool
}

type GetInput struct {
	ExternalKey string
}
type GetOutput struct {
	ExternalKey    string
	DocumentNumber string
}
