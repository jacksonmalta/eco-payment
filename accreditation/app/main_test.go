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
	if input.DocumentNumber == "11111111112" {
		return nil, errors.New("insert error")
	}

	if input.DocumentNumber == "11111111113" {
		return &InsertOutput{
			AlreadyExists: true,
		}, nil
	}

	return &InsertOutput{
		AlreadyExists: false,
	}, nil
}
func (r repositoryMock) GetWithContext(ctx context.Context, input *GetInput) (*GetOutput, error) {
	if r.v == "1" {
		return nil, errors.New("get error")
	}
	if r.v != "" {
		val, err := json.Marshal(input)
		assert.Nil(r.t, err)
		assert.Equal(r.t, r.v, string(val))
		return &GetOutput{
			DocumentNumber: "1",
			ExternalKey:    "2",
		}, nil
	}
	return nil, nil
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

func TestAccreditation_CreateAccountWhenCPF(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"ExternalKey\":\"123\",\"DocumentNumber\":\"11111111111\"}", t)
	a := New(r, l)
	i := &CreateAccountInput{
		DocumentNumber: "11111111111",
		ExternalKey:    "123",
	}
	res, err := a.CreateAccountWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Error\":false,\"Code\":\"\",\"Detail\":\"\"}", string(validate))
}

func TestAccreditation_CreateAccountWhenCNPJ(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"ExternalKey\":\"123\",\"DocumentNumber\":\"11111111111111\"}", t)
	a := New(r, l)
	i := &CreateAccountInput{
		DocumentNumber: "11111111111111",
		ExternalKey:    "123",
	}
	res, err := a.CreateAccountWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Error\":false,\"Code\":\"\",\"Detail\":\"\"}", string(validate))
}

func TestAccreditation_NotCreateAccountWhenInvalidDocument(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("", t)
	a := New(r, l)
	i := &CreateAccountInput{
		DocumentNumber: "111111111111-11",
		ExternalKey:    "123",
	}
	res, err := a.CreateAccountWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Error\":true,\"Code\":\"document-invalid\",\"Detail\":\"invalid document number\"}", string(validate))
}

func TestAccreditation_NotCreateAccountWhenInsertError(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"ExternalKey\":\"123\",\"DocumentNumber\":\"11111111112\"}", t)
	a := New(r, l)
	i := &CreateAccountInput{
		DocumentNumber: "11111111112",
		ExternalKey:    "123",
	}
	res, err := a.CreateAccountWithContext(context.Background(), i)
	assert.Equal(t, "insert error", err.Error())
	assert.Nil(t, res)
}

func TestAccreditation_NotCreateAccountWhenItemAlreadyExists(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"ExternalKey\":\"123\",\"DocumentNumber\":\"11111111113\"}", t)
	a := New(r, l)
	i := &CreateAccountInput{
		DocumentNumber: "11111111113",
		ExternalKey:    "123",
	}
	res, err := a.CreateAccountWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Error\":true,\"Code\":\"item-already-exists\",\"Detail\":\"item already exists\"}", string(validate))
}

func TestAccreditation_GetAccount(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("{\"ExternalKey\":\"1\"}", t)
	a := New(r, l)
	i := &GetAccountInput{
		ExternalKey: "1",
	}
	res, err := a.GetAccountWithContext(context.Background(), i)
	assert.Nil(t, err)
	validate, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"DocumentNumber\":\"1\",\"ExternalKey\":\"2\"}", string(validate))
}

func TestAccreditation_NotGetAccountWhenGetError(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("1", t)
	a := New(r, l)
	i := &GetAccountInput{
		ExternalKey: "1",
	}
	res, err := a.GetAccountWithContext(context.Background(), i)
	assert.Nil(t, res)
	assert.Equal(t, "get error", err.Error())
}

func TestAccreditation_NotGetAccountWhenGetNotFound(t *testing.T) {
	l := newLogMock()
	r := newRepositoryMock("", t)
	a := New(r, l)
	i := &GetAccountInput{
		ExternalKey: "1",
	}
	res, err := a.GetAccountWithContext(context.Background(), i)
	assert.Nil(t, err)
	assert.Nil(t, res)
}
