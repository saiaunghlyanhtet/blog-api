package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/saiaunghlyanhtet/blog-api/pkg/handlers"
	"github.com/saiaunghlyanhtet/blog-api/pkg/post"
)

var dynamodbClient dynamodbiface.DynamoDBAPI

const tableName = "blog-db"

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return
	}

	dynamodbClient = dynamodb.New(awsSession)
	post.InitializeS3SessionAndBucket(awsSession, "blog-api-s3-bucket")
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		if req.PathParameters["id"] != "" {
			return handlers.GetPostById(req, tableName, dynamodbClient)
		} else {
			return handlers.GetAllPostsOverview(req, tableName, dynamodbClient)
		}
	case "POST":
		return handlers.CreatePost(req, tableName, dynamodbClient)
	case "DELETE":
		return handlers.DeletePost(req, tableName, dynamodbClient)
	default:
		return handlers.UnhandledMethod()
	}
}
