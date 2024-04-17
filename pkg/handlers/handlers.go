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

func CreatePost(req events.APIGatewayProxyRequest, tableName string, dynamoDBClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	err := post.CreatePost(dynamoDBClient, tableName, req)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusCreated, "Post created successfully.")
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

func GetPostById(req events.APIGatewayProxyRequest, tableName string, dynamoDBClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	createdDate := req.QueryStringParameters["createdDate"]
	result, err := post.GetPostById(dynamoDBClient, tableName, id, createdDate)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, result)
}

func DeletePost(req events.APIGatewayProxyRequest, tableName string, dynamoDBClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	err := post.DeletePost(dynamoDBClient, tableName, id)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, "Delete the post with ID: "+id)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, "Method Not Allowed.")
}
