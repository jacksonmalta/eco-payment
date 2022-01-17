package routes

import (
	"balance/app"
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

func (r *accreditationMock) SettlementWithContext(ctx context.Context, input *app.SettlementInput) (*app.SettlementOutput, error) {
	vt, err := json.Marshal(input)
	assert.Nil(r.t, err)
	assert.Equal(r.t, r.v, string(vt))

	if input.AccountKey == "12345" {
		return nil, errors.New("settlement error")
	}

	if input.AccountKey == "1234567" {
		return &app.SettlementOutput{
			Error:  true,
			Code:   "item-already-exists",
			Detail: "test2",
		}, nil
	}

	return nil, nil
}
func newAccreditationMock(v string, t *testing.T) app.Balance {
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

func TestRoutes_Settlement(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"external_key\": \"1234\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("{\"AccountKey\":\"123\",\"ExternalKey\":\"1234\",\"OperationType\":\"credit\",\"Amount\":1000}", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestRoutes_NotTestRoutes_SettlementWhenInvalidPayload(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader(""))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"invalid payload\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenAccountKeyEmpty(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"\", \"external_key\": \"1234\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"account_key is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenAccountKeyMissing(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"external_key\": \"1234\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"account_key is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenExternalKeyEmpty(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"external_key\": \"\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"external_key is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenExternalKeyMissing(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"external_key is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenOperationTypeEmpty(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"external_key\": \"1234\", \"operation_type\": \"\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"operation_type is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenOperationTypeMissing(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"external_key\": \"1234\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"operation_type is missing or null\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenAmountZero(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"external_key\": \"1234\", \"operation_type\": \"credit\", \"amount\": 0}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"amount is missing or 0\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotSettlementWhenAmountMissing(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"123\", \"external_key\": \"1234\", \"operation_type\": \"credit\"}"))
	accreditation := newAccreditationMock("", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	expected := "{\"error\":{\"type\":\"invalid_request\",\"category\":\"bad_request\",\"message\":\"amount is missing or 0\"}}"
	assert.Equal(t, expected, string(validate))
}

func TestRoutes_NotBalanceWhenSettlementError(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"12345\", \"external_key\": \"1234\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("{\"AccountKey\":\"12345\",\"ExternalKey\":\"1234\",\"OperationType\":\"credit\",\"Amount\":1000}", t)
	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
	assert.Nil(t, res)
	assert.Equal(t, "settlement error", err.Error())
}

func TestRoutes_NotBalanceWhenSettlementItemAlreadyExists(t *testing.T) {
	l := newLogMock()
	rc := io.NopCloser(strings.NewReader("{\"account_key\": \"1234567\", \"external_key\": \"1234\", \"operation_type\": \"credit\", \"amount\": 1000}"))
	accreditation := newAccreditationMock("{\"AccountKey\":\"1234567\",\"ExternalKey\":\"1234\",\"OperationType\":\"credit\",\"Amount\":1000}", t)

	res, err := balanceWithContext(context.Background(), rc, l, accreditation)
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

func TestRoutes_IntValueWhenNilValue(t *testing.T) {
	r := intValue(nil)
	assert.Equal(t, 0, r)
}
