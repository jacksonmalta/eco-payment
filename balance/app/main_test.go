package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type repositoryMock struct {
	v string
	t *testing.T
}

func (r repositoryMock) InsertWithContext(ctx context.Context, input *InsertInput) (*InsertOutput, error) {
	v, err := json.Marshal(input)
	assert.Nil(r.t, err)
	assert.Equal(r.t, r.v, string(v))
	if input.AccountKey == "11111111112" {
		return nil, errors.New("insert error")
	}

	if input.AccountKey == "11111111113" {
		return &InsertOutput{
			AlreadyExists: true,
		}, nil
	}

	return &InsertOutput{
		AlreadyExists: false,
	}, nil
}
func newRepositoryMock(v string, t *testing.T) Persistence {
	return &repositoryMock{
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

func TestAccreditation_Settlement(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"AccountKey\":\"11111111111\",\"ExternalKey\":\"123\",\"OperatiionType\":\"test\",\"Amount\":1000}", t)
	a := New(r, l)
	i := &SettlementInput{
		AccountKey:    "11111111111",
		ExternalKey:   "123",
		OperationType: "test",
		Amount:        1000,
	}
	res, err := a.SettlementWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Error\":false,\"Code\":\"\",\"Detail\":\"\"}", string(validate))
}

func TestAccreditation_NotSettlementWhenInsertError(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"AccountKey\":\"11111111112\",\"ExternalKey\":\"123\",\"OperatiionType\":\"test\",\"Amount\":1000}", t)
	a := New(r, l)
	i := &SettlementInput{
		AccountKey:    "11111111112",
		ExternalKey:   "123",
		OperationType: "test",
		Amount:        1000,
	}
	res, err := a.SettlementWithContext(context.Background(), i)
	assert.Equal(t, "insert error", err.Error())
	assert.Nil(t, res)
}

func TestAccreditation_NotSettlementWhenItemAlreadyExists(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"AccountKey\":\"11111111113\",\"ExternalKey\":\"123\",\"OperatiionType\":\"test\",\"Amount\":1000}", t)
	a := New(r, l)
	i := &SettlementInput{
		AccountKey:    "11111111113",
		ExternalKey:   "123",
		OperationType: "test",
		Amount:        1000,
	}
	res, err := a.SettlementWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Error\":true,\"Code\":\"item-already-exists\",\"Detail\":\"item already exists\"}", string(validate))
}
