package routes

import (
	"balance/app"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	BadRequest     = "bad_request"
	Conflict       = "conflict"
	InvalidRequest = "invalid_request"
)

func stringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func intValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

type BalanceRequest struct {
	AccountKey    *string `json:"account_key,omitempty"`
	ExternalKey   *string `json:"external_key,omitempty"`
	OperationType *string `json:"operation_type,omitempty"`
	Amount        *int    `json:"amount,omitempty"`
}

type BalanceError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type,omitempty"`
	Category   string `json:"category,omitempty"`
	Message    string `json:"message,omitempty"`
}

type BalanceErrorResponse struct {
	Error *BalanceError `json:"error,omitempty"`
}

func responseBuild(msg string, statusCode int, category string) *BalanceErrorResponse {
	et := &BalanceError{
		StatusCode: statusCode,
		Type:       InvalidRequest,
		Category:   category,
		Message:    msg,
	}
	ae := &BalanceErrorResponse{
		Error: et,
	}
	return ae
}

func buildBalanceRequest(a []byte) (*BalanceRequest, *BalanceErrorResponse) {
	va := &BalanceRequest{}

	err := json.Unmarshal(a, &va)
	if err != nil {
		return nil, responseBuild("invalid payload", http.StatusBadRequest, BadRequest)
	}

	if va.AccountKey == nil || stringValue(va.AccountKey) == "" {
		return nil, responseBuild("account_key is missing or null", http.StatusBadRequest, BadRequest)
	}

	if va.ExternalKey == nil || stringValue(va.ExternalKey) == "" {
		return nil, responseBuild("external_key is missing or null", http.StatusBadRequest, BadRequest)
	}

	if va.OperationType == nil || stringValue(va.OperationType) == "" {
		return nil, responseBuild("operation_type is missing or null", http.StatusBadRequest, BadRequest)
	}

	if va.Amount == nil || intValue(va.Amount) == 0 {
		return nil, responseBuild("amount is missing or 0", http.StatusBadRequest, BadRequest)
	}

	return va, nil
}

func balanceWithContext(ctx context.Context, body io.ReadCloser, log Logger, a app.Balance) (*BalanceErrorResponse, error) {
	defer body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	b := buf.Bytes()

	request, errorResponse := buildBalanceRequest(b)

	if errorResponse != nil {
		return errorResponse, nil
	}

	i := &app.SettlementInput{
		AccountKey:    stringValue(request.AccountKey),
		ExternalKey:   stringValue(request.ExternalKey),
		OperationType: stringValue(request.OperationType),
		Amount:        intValue(request.Amount),
	}

	res, err := a.SettlementWithContext(ctx, i)

	if err != nil {
		return nil, err
	}

	if res != nil && res.Error && res.Code == "item-already-exists" {
		return responseBuild(res.Detail, http.StatusConflict, Conflict), nil
	}

	return nil, nil
}
