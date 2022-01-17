package authorizer

import (
	"context"
	"debit/app"
	"fmt"
	"net/http"
)

type accreditation struct {
	log         Logger
	config      *Config
	httpService Http
}

type Http interface {
	GetWithContext(ctx context.Context, url string) ([]byte, int, error)
}

func (a *accreditation) AuthorizeWithContext(ctx context.Context, input *app.AuthorizeInput) (*app.AuthorizeOutput, error) {
	_, statusCode, err := a.httpService.GetWithContext(ctx, a.config.Url+input.AccountKey)
	if err != nil {
		a.log.Error(fmt.Sprintf("http get error %s", err.Error()))
		return nil, err
	}

	if statusCode == http.StatusNotFound {
		return nil, nil
	}

	if statusCode >= http.StatusOK {
		return &app.AuthorizeOutput{
			HasError: false,
		}, nil
	}

	return &app.AuthorizeOutput{
		HasError: true,
	}, nil
}

func New(log Logger, config *Config, httpService Http) app.Authorizer {
	return &accreditation{
		log:         log,
		config:      config,
		httpService: httpService,
	}
}
