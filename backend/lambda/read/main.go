package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Secret struct {
	SecretID   	   string  `json:"SecretID,omitempty"`
	Message    	   string  `json:"Message"`
	SecretKey  	   string  `json:"SecretKey"`
	ActiveDuration string  `json:"ActiveDuration"`
	ExpirationTime int64   `json:"ExpirationTime,omitempty"`
}

var (
	headers = map[string]string{"Content-type":"application/json"}
)

func epochValue(duration time.Duration) int64 {
	curTime := time.Now()
	log.Printf("Current Time in Unix [%d]", curTime.Unix())
	expTime := curTime.Add(time.Hour * duration)
	epoch := expTime.Unix()
	log.Printf("Calculated EPOCH value is [%d]", epoch)
	return epoch
}

func dynamodbHandler(s Secret, sid string) error {
	sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION"))},
    )
    svc := dynamodb.New(sess)
	dynamoMap, err := dynamodbattribute.MarshalMap(s)

	if err != nil {
        log.Printf("Error occurred while Marshalling DynamoDB data [%s]", err)
		return errors.New("INTERNAL ERROR. PLEASE TRY AGAIN")
    }

	putItem := &dynamodb.PutItemInput{
        Item: dynamoMap,
        TableName: aws.String(os.Getenv("DDB_TABLE")),
    }
	putRes, err := svc.PutItem(putItem)

	if err != nil {
        log.Printf("Error occurred during DynamoDB PUT operation [%s]", err)
		return errors.New("INTERNAL ERROR. PLEASE TRY AGAIN")
    }

    log.Printf("Secret [%s] has been added successfully - [%s]", sid, putRes)

	return nil
}

func lambdaHandler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	var s Secret
	
	sid := strings.ReplaceAll(request.RequestContext.RequestID, "-", "")
	s.SecretID = sid

	err := json.Unmarshal([]byte(request.Body), &s)

	if err != nil {
		log.Printf("Error occurred while Unmarshalling the JSON request [%s]",err)
		return events.LambdaFunctionURLResponse{Body:"INTERNAL ERROR. PLEASE TRY AGAIN", Headers: headers, StatusCode: 200}, nil
	}

	switch s.ActiveDuration {
		case "1h":
			log.Printf("ActiveDuration is [%s]", s.ActiveDuration)
			s.ExpirationTime = epochValue(1)
		case "2h":
			log.Printf("ActiveDuration is [%s]", s.ActiveDuration)
			s.ExpirationTime = epochValue(2)
		case "12h":
			log.Printf("ActiveDuration is [%s]", s.ActiveDuration)
			s.ExpirationTime = epochValue(12)
		case "24h":
			log.Printf("ActiveDuration is [%s]", s.ActiveDuration)
			s.ExpirationTime = epochValue(24)
		default:
			log.Printf("ActiveDuration is [%s] and setting default value.", s.ActiveDuration)
			s.ActiveDuration = "1h"
			s.ExpirationTime = epochValue(1)
	}

	err = dynamodbHandler(s, sid) 

	if err != nil {
		return events.LambdaFunctionURLResponse{Body:err.Error(), Headers: headers, StatusCode: 200}, nil	
	}

	return events.LambdaFunctionURLResponse{Body:sid, Headers: headers, StatusCode: 200}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}