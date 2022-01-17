#!bin/bash

export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=bar

aws --endpoint-url=http://localhost:4566 dynamodb create-table \
    --table-name account \
    --attribute-definitions \
        AttributeName=ExternalKey,AttributeType=S \
    --key-schema \
        AttributeName=ExternalKey,KeyType=HASH \
    --billing-mode \
        PAY_PER_REQUEST \


