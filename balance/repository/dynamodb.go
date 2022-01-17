package repository

import (
	"balance/app"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

type Dynamodb interface {
	PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

type db struct {
	dynamodbService Dynamodb
	log             Logger
	config          Config
}

func (d *db) InsertWithContext(ctx context.Context, input *app.InsertInput) (*app.InsertOutput, error) {
	putItemInput := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"AccountKey": {
				S: aws.String(input.AccountKey),
			},
			"ExternalKey": {
				S: aws.String(input.ExternalKey),
			},
			"OperationType": {
				S: aws.String(input.OperatiionType),
			},
			"Amount": {
				N: aws.String(strconv.Itoa(input.Amount)),
			},
		},
		TableName:           aws.String(d.config.TableName),
		ConditionExpression: aws.String("attribute_not_exists(AccountKey) AND attribute_not_exists(ExternalKey)"),
	}
	d.log.Info(fmt.Sprintf("Dynamodb input item %v", input))
	_, err := d.dynamodbService.PutItemWithContext(ctx, putItemInput)
	if err != nil {
		if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
			d.log.Info(fmt.Sprintf("%s", ae.Code()))
			return &app.InsertOutput{
				AlreadyExists: true,
			}, nil
		}
		d.log.Error(fmt.Sprintf("Error %s", err.Error()))
		return nil, err
	}

	return &app.InsertOutput{
		AlreadyExists: false,
	}, nil
}

func NewDynamodb(d Dynamodb, log Logger, config Config) app.Persistence {
	return &db{
		dynamodbService: d,
		log:             log,
		config:          config,
	}
}
