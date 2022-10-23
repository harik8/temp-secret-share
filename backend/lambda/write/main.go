package main

import (
	"context"
	"log"
	"strconv"
	"time"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	headers = map[string]string{"Content-type":"application/json"}
)

type Secret struct {
	SecretID   	   string  `json:"SecretID,omitempty"`
	Message    	   string  `json:"Message"`
	SecretKey  	   string  `json:"SecretKey"`
	ActiveDuration string  `json:"ActiveDuration"`
	ExpirationTime int64   `json:"ExpirationTime,omitempty"`
}

func currentEpochValue() int64 {
	curTime := time.Now()
	log.Printf("Current Time in Epoch [%d]", curTime.Unix())
	return curTime.Unix()
}

func delSecret(sid string, expTime int64, svc *dynamodb.DynamoDB) error {

	input := &dynamodb.DeleteItemInput{
        Key : map[string]*dynamodb.AttributeValue{
			"SecretID": {
				S: aws.String(sid),
			},
			"ExpirationTime": {
				N: aws.String(strconv.Itoa(int(expTime))),
			},
    	},
		TableName : aws.String(os.Getenv("DDB_TABLE")),
	}

	_, err := svc.DeleteItem(input)

	return err
}

func lambdaHandler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	var s Secret

	log.Printf("QueryStringParameters - SecretID [%s]", request.QueryStringParameters["SecretID"])

	sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION"))},
    )
    svc := dynamodb.New(sess)

	pkCond  := expression.Key("SecretID").Equal(expression.Value(request.QueryStringParameters["SecretID"]))
	skCond  := expression.Key("ExpirationTime").GreaterThan(expression.Value(currentEpochValue()))
	filCond := expression.Name("SecretKey").Equal(expression.Value(request.QueryStringParameters["SecretKey"]))

	expr, err := expression.NewBuilder().
		WithKeyCondition(pkCond.And(skCond)).
		WithFilter(filCond).
		Build()

	if err != nil {
		log.Printf("Error occurred while building expression [%s]", err)
		return events.LambdaFunctionURLResponse{Body:"INTERNAL ERROR. PLEASE TRY AGAIN.", Headers: headers, StatusCode: 400}, nil
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames  : expr.Names(),
		ExpressionAttributeValues : expr.Values(),
		KeyConditionExpression	  : expr.KeyCondition(),
		FilterExpression		  : expr.Filter(),
		TableName				  : aws.String("secrets-share"),
	}

	result, err := svc.Query(input)

	if err != nil {
		log.Printf("Error occurred while Querying [%s]", err)
		return events.LambdaFunctionURLResponse{Body:"EITHER SECRET IS EXPIRED, ALREADY READ OR INCORRECT SECRET KEY", Headers: headers, StatusCode: 200}, nil
	}

	if *result.Count == 0 {
		log.Printf("Secret not found")
		return events.LambdaFunctionURLResponse{Body:"EITHER SECRET IS EXPIRED, ALREADY READ OR INCORRECT SECRET KEY", Headers: headers, StatusCode: 200}, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &s)

	if err != nil {
		log.Printf("Error occurred while Unmarshalling [%s]", err)
		return events.LambdaFunctionURLResponse{Body:"INTERNAL ERROR. PLEASE TRY AGAIN.", Headers: headers, StatusCode: 200}, nil
	}

	log.Printf("The Secret ID [%s] has been retrieved", s.SecretID)

	err = delSecret(s.SecretID, s.ExpirationTime, svc)

	if err != nil {
		log.Printf("Error occurred while Deleting the Secret [%s]. Error [%s]", s.SecretID, err)
		return events.LambdaFunctionURLResponse{Body:s.Message, Headers: headers, StatusCode: 200}, nil
	}

	return events.LambdaFunctionURLResponse{Body:s.Message, Headers: headers, StatusCode: 200}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}