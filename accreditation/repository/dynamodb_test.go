package repository

import (
	"accreditation/app"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"testing"
)

type serviceMock struct {
	v string
	t *testing.T
}

type ea struct{}

func (e ea) Error() string {
	return ""
}

func (e ea) Code() string {
	return "ConditionalCheckFailedException"
}
func (e ea) Message() string {
	return ""
}
func (e ea) OrigErr() error {
	return nil
}
func (e ea) StatusCode() int {
	return 0
}
func (e ea) RequestID() string {
	return ""
}

func ErrorAws() awserr.RequestFailure {
	return &ea{}
}

func (s serviceMock) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if s.v != "" {
		if s.v == "1" {
			return nil, ErrorAws()
		}
		v, err := json.Marshal(input)
		assert.Nil(s.t, err)
		assert.Equal(s.t, s.v, string(v))
		return nil, nil
	}
	return nil, errors.New("db error")
}
func (s serviceMock) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	if s.v == "1" {
		return nil, errors.New("get error")
	}
	if s.v != "" {
		val, err := json.Marshal(input)
		assert.Nil(s.t, err)
		assert.Equal(s.t, s.v, string(val))
		attributeValue := make(map[string]*dynamodb.AttributeValue)
		attributeValue["ExternalKey"] = &dynamodb.AttributeValue{
			S: aws.String("1"),
		}
		attributeValue["DocumentNumber"] = &dynamodb.AttributeValue{
			S: aws.String("2"),
		}
		return &dynamodb.GetItemOutput{
			Item: attributeValue,
		}, nil
	}
	return nil, nil
}
func newServiceMock(v string, t *testing.T) Dynamodb {
	return &serviceMock{
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

func TestDb_Insert(t *testing.T) {
	l := newLogMock()
	s := newServiceMock("{\"ConditionExpression\":\"attribute_not_exists(ExternalKey)\",\"ConditionalOperator\":null,\"Expected\":null,\"ExpressionAttributeNames\":null,\"ExpressionAttributeValues\":null,\"Item\":{\"DocumentNumber\":{\"B\":null,\"BOOL\":null,\"BS\":null,\"L\":null,\"M\":null,\"N\":null,\"NS\":null,\"NULL\":null,\"S\":\"1\",\"SS\":null},\"ExternalKey\":{\"B\":null,\"BOOL\":null,\"BS\":null,\"L\":null,\"M\":null,\"N\":null,\"NS\":null,\"NULL\":null,\"S\":\"2\",\"SS\":null}},\"ReturnConsumedCapacity\":null,\"ReturnItemCollectionMetrics\":null,\"ReturnValues\":null,\"TableName\":\"account\"}", t)
	c := Config{
		TableName: "account",
	}
	d := NewDynamodb(s, l, c)
	i := &app.InsertInput{
		DocumentNumber: "1",
		ExternalKey:    "2",
	}
	res, err := d.InsertWithContext(context.Background(), i)
	assert.Nil(t, err)
	b, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"AlreadyExists\":false}", string(b))
}

func TestDb_NotInsert(t *testing.T) {
	l := newLogMock()
	s := newServiceMock("", t)
	c := Config{
		TableName: "account",
	}
	d := NewDynamodb(s, l, c)
	i := &app.InsertInput{
		DocumentNumber: "1",
		ExternalKey:    "2",
	}
	res, err := d.InsertWithContext(context.Background(), i)
	assert.Nil(t, res)
	assert.Equal(t, "db error", err.Error())
}

func TestDb_NotInsertWhenExternalKeyHasExists(t *testing.T) {
	l := newLogMock()
	s := newServiceMock("1", t)
	c := Config{
		TableName: "account",
	}
	d := NewDynamodb(s, l, c)
	i := &app.InsertInput{
		DocumentNumber: "1",
		ExternalKey:    "2",
	}
	res, err := d.InsertWithContext(context.Background(), i)
	assert.Nil(t, err)
	b, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"AlreadyExists\":true}", string(b))
}

func TestDb_Get(t *testing.T) {
	l := newLogMock()
	s := newServiceMock("{\"AttributesToGet\":null,\"ConsistentRead\":null,\"ExpressionAttributeNames\":null,\"Key\":{\"ExternalKey\":{\"B\":null,\"BOOL\":null,\"BS\":null,\"L\":null,\"M\":null,\"N\":null,\"NS\":null,\"NULL\":null,\"S\":\"1\",\"SS\":null}},\"ProjectionExpression\":null,\"ReturnConsumedCapacity\":null,\"TableName\":\"account\"}", t)
	c := Config{
		TableName: "account",
	}
	d := NewDynamodb(s, l, c)
	i := &app.GetInput{
		ExternalKey: "1",
	}
	res, err := d.GetWithContext(context.Background(), i)
	assert.Nil(t, err)
	b, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.Equal(t, "{\"ExternalKey\":\"1\",\"DocumentNumber\":\"2\"}", string(b))
}

func TestDb_NotGetWhenGetItemError(t *testing.T) {
	l := newLogMock()
	s := newServiceMock("1", t)
	c := Config{
		TableName: "account",
	}
	d := NewDynamodb(s, l, c)
	i := &app.GetInput{
		ExternalKey: "1",
	}
	res, err := d.GetWithContext(context.Background(), i)
	assert.Nil(t, res)
	assert.Equal(t, "get error", err.Error())
}

func TestDb_NotGetWhenGetItemNotFound(t *testing.T) {
	l := newLogMock()
	s := newServiceMock("", t)
	c := Config{
		TableName: "account",
	}
	d := NewDynamodb(s, l, c)
	i := &app.GetInput{
		ExternalKey: "1",
	}
	res, err := d.GetWithContext(context.Background(), i)
	assert.Nil(t, res)
	assert.Nil(t, err)
}
