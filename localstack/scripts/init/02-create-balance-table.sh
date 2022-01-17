#!bin/bash

export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=bar

aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name balance \
    --attribute-definitions \
        AttributeName=AccountKey,AttributeType=S \
        AttributeName=ExternalKey,AttributeType=S \
    --key-schema \
        AttributeName=AccountKey,KeyType=HASH \
        AttributeName=ExternalKey,KeyType=RANGE \
    --billing-mode \
        PAY_PER_REQUEST \


