package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"

	_ "go.elastic.co/apm/module/apmlambda/v2"
)

var coldstart = true

func Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("Example function log %s %v", lc.AwsRequestID, coldstart)
	var body struct {
		Sleep bool `json:"sleep"`
	}
	err := json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Hello from go!%s", err),
			Headers: map[string]string{
				"coldstart": strconv.FormatBool(coldstart),
			},
		}, err
	}
	if body.Sleep {
		log.Println("decided to sleep")
		time.Sleep(20 * time.Second)
	} else {
		log.Println("decided to NOT sleep")
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello from go!%s", lc.AwsRequestID),
		Headers: map[string]string{
			"coldstart": strconv.FormatBool(coldstart),
		},
	}
	coldstart = false
	return response, nil
}

func main() {
	lambda.Start(Handle)
}
