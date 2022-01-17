package app

import "context"

type Authorizer interface {
	AuthorizeWithContext(ctx context.Context, input *AuthorizeInput) (*AuthorizeOutput, error)
}

type AuthorizeInput struct {
	AccountKey string
}
type AuthorizeOutput struct {
	HasError bool
}
