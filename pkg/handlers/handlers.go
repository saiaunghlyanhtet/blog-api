package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/saiaunghlyanhtet/blog-api/pkg/post"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetAllPostsOverview(req events.APIGatewayProxyRequest, tableName string, dynamoDBClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := post.GetAllPostsOverview(dynamoDBClient, tableName)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, result)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, "Method Not Allowed.")
}
