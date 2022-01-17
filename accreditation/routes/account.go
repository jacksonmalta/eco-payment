package routes

import (
	"accreditation/app"
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

type AccountRequest struct {
	DocumentNumber *string `json:"document_number,omitempty"`
	ExternalKey    *string `json:"external_key,omitempty"`
}

type AccountError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type,omitempty"`
	Category   string `json:"category,omitempty"`
	Message    string `json:"message,omitempty"`
}

type AccountErrorResponse struct {
	Error *AccountError `json:"error,omitempty"`
}

type AccountGetResponse struct {
	DocumentNumber string `json:"document_number,omitempty"`
	ExternalKey    string `json:"external_key,omitempty"`
}

func responseBuild(msg string, statusCode int, category string) *AccountErrorResponse {
	et := &AccountError{
		StatusCode: statusCode,
		Type:       InvalidRequest,
		Category:   category,
		Message:    msg,
	}
	ae := &AccountErrorResponse{
		Error: et,
	}
	return ae
}

func buildAccountRequest(a []byte) (*AccountRequest, *AccountErrorResponse) {
	va := &AccountRequest{}

	err := json.Unmarshal(a, &va)
	if err != nil {
		return nil, responseBuild("invalid payload", http.StatusBadRequest, BadRequest)
	}

	if va.DocumentNumber == nil || stringValue(va.DocumentNumber) == "" {
		return nil, responseBuild("document_number is missing or null", http.StatusBadRequest, BadRequest)
	}

	if va.ExternalKey == nil || stringValue(va.ExternalKey) == "" {
		return nil, responseBuild("external_key is missing or null", http.StatusBadRequest, BadRequest)
	}

	return va, nil
}

func createAccountWithContext(ctx context.Context, body io.ReadCloser, log Logger, a app.Accreditation) (*AccountErrorResponse, error) {
	defer body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	b := buf.Bytes()

	accountRequest, accountErrorResponse := buildAccountRequest(b)

	if accountErrorResponse != nil {
		return accountErrorResponse, nil
	}

	i := &app.CreateAccountInput{
		DocumentNumber: stringValue(accountRequest.DocumentNumber),
		ExternalKey:    stringValue(accountRequest.ExternalKey),
	}

	res, err := a.CreateAccountWithContext(ctx, i)

	if err != nil {
		return nil, err
	}

	if res != nil && res.Error && res.Code == "document-invalid" {
		return responseBuild(res.Detail, http.StatusBadRequest, BadRequest), nil
	}

	if res != nil && res.Error && res.Code == "item-already-exists" {
		return responseBuild(res.Detail, http.StatusConflict, Conflict), nil
	}

	return nil, nil
}

func getAccountWithContext(ctx context.Context, externalKey string, log Logger, a app.Accreditation) (*AccountGetResponse, error) {
	i := &app.GetAccountInput{
		ExternalKey: externalKey,
	}

	res, err := a.GetAccountWithContext(ctx, i)

	if err != nil {
		return nil, err
	}

	if res != nil {
		return &AccountGetResponse{
			ExternalKey:    res.ExternalKey,
			DocumentNumber: res.DocumentNumber,
		}, nil
	}

	return nil, nil
}
