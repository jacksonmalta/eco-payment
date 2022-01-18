package routes

import (
	"bytes"
	"context"
	"credit/app"
	"encoding/json"
	"io"
	"net/http"
)

const (
	BadRequest     = "bad_request"
	Conflict       = "conflict"
	InvalidRequest = "invalid_request"
	BadGateway     = "bad_gateway"
	NotFound       = "not_found"
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

type TransactionRequest struct {
	AccountKey  *string `json:"account_key,omitempty"`
	ExternalKey *string `json:"external_key,omitempty"`
	Amount      *int    `json:"amount,omitempty"`
}

type TransactionError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type,omitempty"`
	Category   string `json:"category,omitempty"`
	Message    string `json:"message,omitempty"`
}

type TransactionErrorResponse struct {
	Error *TransactionError `json:"error,omitempty"`
}

func responseBuild(msg string, statusCode int, category string) *TransactionErrorResponse {
	et := &TransactionError{
		StatusCode: statusCode,
		Type:       InvalidRequest,
		Category:   category,
		Message:    msg,
	}
	ae := &TransactionErrorResponse{
		Error: et,
	}
	return ae
}

func buildTransactionRequest(a []byte) (*TransactionRequest, *TransactionErrorResponse) {
	va := &TransactionRequest{}

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

	if va.Amount == nil || intValue(va.Amount) == 0 {
		return nil, responseBuild("amount is missing or 0", http.StatusBadRequest, BadRequest)
	}

	return va, nil
}

func transactionWithContext(ctx context.Context, body io.ReadCloser, log Logger, a app.Credit) (*TransactionErrorResponse, error) {
	defer body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	b := buf.Bytes()

	request, errorResponse := buildTransactionRequest(b)

	if errorResponse != nil {
		return errorResponse, nil
	}

	i := &app.TransactionInput{
		AccountKey:  stringValue(request.AccountKey),
		ExternalKey: stringValue(request.ExternalKey),
		Amount:      intValue(request.Amount),
	}

	res, err := a.TransactionWithContext(ctx, i)

	if err != nil {
		return nil, err
	}

	if res != nil && res.Error && res.Code == app.UnauthorizedTransaction {
		return responseBuild(res.Detail, http.StatusBadGateway, BadGateway), nil
	}

	if res != nil && res.Error && res.Code == app.UnauthorizedSettlement {
		return responseBuild(res.Detail, http.StatusBadGateway, BadGateway), nil
	}

	if res != nil && res.Error && res.Code == app.AuthorizerNotFound {
		return responseBuild("Account Key not found", http.StatusNotFound, NotFound), nil
	}

	if res != nil && res.Error && res.Code == app.SettlementFailed {
		return responseBuild(res.Detail, http.StatusConflict, Conflict), nil
	}

	return nil, nil
}
