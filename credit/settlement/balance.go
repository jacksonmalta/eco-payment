package settlement

import (
	"context"
	"credit/app"
	"encoding/json"
	"fmt"
	"net/http"
)

type balance struct {
	log         Logger
	config      *Config
	httpService Http
}

type Http interface {
	PostWithContext(ctx context.Context, url string, payload []byte) ([]byte, int, error)
}

type BalanceRequestPayload struct {
	AccountKey    string `json:"account_key,omitempty"`
	ExternalKey   string `json:"external_key,omitempty"`
	OperationType string `json:"operation_type,omitempty"`
	Amount        int    `json:"amount,omitempty"`
}

type BalanceError struct {
	Type     string `json:"type,omitempty"`
	Category string `json:"category,omitempty"`
	Message  string `json:"message,omitempty"`
}

type BalanceResponseError struct {
	Error *BalanceError `json:"error,omitempty"`
}

func (b *balance) SettleWithContext(ctx context.Context, input *app.SettleInput) (*app.SettleOutput, error) {
	payload := &BalanceRequestPayload{
		AccountKey:    input.AccountKey,
		ExternalKey:   input.ExternalKey,
		OperationType: input.OperationType,
		Amount:        input.Amount,
	}
	pb, err := json.Marshal(payload)
	res, statusCode, err := b.httpService.PostWithContext(ctx, b.config.Url, pb)
	if err != nil {
		b.log.Error(fmt.Sprintf("http post error %s", err.Error()))
		return nil, err
	}

	if statusCode == http.StatusBadRequest || statusCode == http.StatusConflict {
		be := &BalanceResponseError{}
		err := json.Unmarshal(res, be)
		if err != nil {
			b.log.Error(fmt.Sprintf("http post error %s", err.Error()))
			return nil, err
		}
		return &app.SettleOutput{
			HasIntermitance: false,
			Error:           true,
			Code:            be.Error.Type,
			Detail:          be.Error.Message,
		}, nil
	}

	if statusCode >= http.StatusCreated {
		return &app.SettleOutput{
			HasIntermitance: false,
			Error:           false,
		}, nil
	}

	return &app.SettleOutput{
		HasIntermitance: true,
	}, nil
}

func New(log Logger, config *Config, httpService Http) app.Settlement {
	return &balance{
		log:         log,
		config:      config,
		httpService: httpService,
	}
}
