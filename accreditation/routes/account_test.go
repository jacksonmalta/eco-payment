package routes

import (
	"accreditation/app"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

type accreditationMock struct {
	v string
	t *testing.T
}

func (r *accreditationMock) CreateAccountWithContext(ctx context.Context, input *app.CreateAccountInput) (*app.CreateAccountOutput, error) {
	vt, err := json.Marshal(input)
	assert.Nil(r.t, err)
	assert.Equal(r.t, r.v, string(vt))

	if input.DocumentNumber == "12345" {
		return nil, errors.New("account error")
	}

	if input.DocumentNumber == "123456" {
		return &app.CreateAccountOutput{
			Error:  true,
			Code:   "document-invalid",
			Detail: "test",
		}, nil
	}

	if input.DocumentNumber == "1234567" {
		return &app.CreateAccountOutput{
			Error:  true,
			Code:   "item-already-exists",
			Detail: "test2",
		}, nil
	}

	return nil, nil
}
func (r *accreditationMock) GetAccountWithContext(ctx context.Context, input *app.GetAccountInput) (*app.GetAccountOutput, error) {
	if r.v == "1" {
		return nil, errors.New("get error")
	}
	if r.v != "" {
		val, err := json.Marshal(input)
		assert.Nil(r.t, err)
		assert.Equal(r.t, r.v, string(val))
		return &app.GetAccountOutput{
			ExternalKey:    "1",
			DocumentNumber: "123",
		}, nil
	}
	return nil, nil
}
func newAccreditationMock(v string, t *testing.T) app.Accreditation {
	return &accreditationMock{
		v: v,
		t: t,
	}
}

type log struct{}

func (l log) Info(msg string)  {}
func (l log) Error(msg string) {}
func newLogMock() Logger {
	return &log{}
}

func TestRoutes_CreateAccount(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"123\", \"external_key\": \"1234\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestRoutes_NotCreateAccountWhenInvalidPayload(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader(""))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"invalid payload\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotCreateAccountWhenDocumentNumberEmpty(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"\", \"external_key\": \"1234\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"document_number is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotCreateAccountWhenDocumentNumberMissing(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"external_key\": \"1234\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"document_number is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotCreateAccountWheExternalKeyEmpty(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"123\", \"external_key\": \"\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"external_key is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotCreateAccountWheExternalKeyMissing(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"123\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"external_key is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotCreateAccountWheAccreditationError(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"12345\", \"external_key\": \"1234\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"12345\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, res)
	assert.Equal(t, "account error", err.Error())
}

func TestRoutes_NotCreateAccountWhenAccreditationInvalidInput(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"123456\", \"external_key\": \"1234\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"123456\",\"ExternalKey\":\"1234\"}", t)
	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"test\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotCreateAccountWhenAccreditationItemAlreadyExists(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"document_number\": \"1234567\", \"external_key\": \"1234\"}"))
	accreditation := newAccreditationMock("{\"DocumentNumber\":\"1234567\",\"ExternalKey\":\"1234\"}", t)

	res, err := createAccountWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"conflict\",\"message\":\"test2\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_StringValueWhenNilValue(t *testing.T) {
	r := stringValue(nil)
	assert.Equal(t, "", r)
}

func TestRoutes_GetAccount(t *testing.T) {
	l := newLogMock()
	accreditation := newAccreditationMock("{\"ExternalKey\":\"1\"}", t)
	res, err := getAccountWithContext(context.Background(), "1", l, accreditation)
	assert.Nil(t, err)
	val, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"document_number\":\"123\",\"external_key\":\"1\"}", string(val))
}

func TestRoutes_NotGetAccountWhenError(t *testing.T) {
	l := newLogMock()
	accreditation := newAccreditationMock("1", t)
	res, err := getAccountWithContext(context.Background(), "1", l, accreditation)
	assert.Nil(t, res)
	assert.Equal(t, "get error", err.Error())
}

func TestRoutes_NotGetAccountWhenNotFound(t *testing.T) {
	l := newLogMock()
	accreditation := newAccreditationMock("", t)
	res, err := getAccountWithContext(context.Background(), "1", l, accreditation)
	assert.Nil(t, res)
	assert.Nil(t, err)
}
