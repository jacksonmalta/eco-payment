package app

import "context"

type Accreditation interface {
	CreateAccountWithContext(ctx context.Context, input *CreateAccountInput) (*CreateAccountOutput, error)
	GetAccountWithContext(ctx context.Context, input *GetAccountInput) (*GetAccountOutput, error)
}

type CreateAccountInput struct {
	DocumentNumber string
	ExternalKey    string
}

type CreateAccountOutput struct {
	Error  bool
	Code   string
	Detail string
}

type GetAccountInput struct {
	ExternalKey string
}

type GetAccountOutput struct {
	DocumentNumber string
	ExternalKey    string
}
