package services

import (
	"balance/repository"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type db struct {
	svc *dynamodb.DynamoDB
}

func (d *db) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return d.svc.PutItemWithContext(ctx, input)
}

func NewDynamodb() repository.Dynamodb {
	mySession := session.Must(session.NewSession())
	svc := dynamodb.New(mySession, aws.NewConfig().WithRegion("us-east-1"), aws.NewConfig().WithEndpoint("http://localstack:4566"))
	return &db{
		svc: svc,
	}
}
