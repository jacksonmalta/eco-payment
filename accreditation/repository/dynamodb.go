package repository

import (
	"accreditation/app"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Dynamodb interface {
	PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

type db struct {
	dynamodbService Dynamodb
	log             Logger
	config          Config
}

func (d *db) InsertWithContext(ctx context.Context, input *app.InsertInput) (*app.InsertOutput, error) {
	putItemInput := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ExternalKey": {
				S: aws.String(input.ExternalKey),
			},
			"DocumentNumber": {
				S: aws.String(input.DocumentNumber),
			},
		},
		TableName:           aws.String(d.config.TableName),
		ConditionExpression: aws.String("attribute_not_exists(ExternalKey)"),
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

func (d *db) GetWithContext(ctx context.Context, input *app.GetInput) (*app.GetOutput, error) {
	attributeValue := make(map[string]*dynamodb.AttributeValue)
	attributeValue["ExternalKey"] = &dynamodb.AttributeValue{
		S: aws.String(input.ExternalKey),
	}
	i := &dynamodb.GetItemInput{
		Key:       attributeValue,
		TableName: aws.String(d.config.TableName),
	}
	getItemOutput, err := d.dynamodbService.GetItemWithContext(ctx, i)
	if err != nil {
		d.log.Error(fmt.Sprintf("Error get item %s", err.Error()))
		return nil, err
	}

	if getItemOutput == nil || getItemOutput.Item == nil || getItemOutput.Item["ExternalKey"] == nil || aws.StringValue(getItemOutput.Item["ExternalKey"].S) == "" {
		return nil, nil
	}

	return &app.GetOutput{
		ExternalKey:    *getItemOutput.Item["ExternalKey"].S,
		DocumentNumber: *getItemOutput.Item["DocumentNumber"].S,
	}, nil
}

func NewDynamodb(d Dynamodb, log Logger, config Config) app.Persistence {
	return &db{
		dynamodbService: d,
		log:             log,
		config:          config,
	}
}
