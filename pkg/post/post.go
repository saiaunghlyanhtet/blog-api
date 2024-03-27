package post

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Post struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Summary string   `json:"summary"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Tags    []string `json:"tags"`
}

// CreatePost creates a new post in the DynamoDB table and puts the image to S3 bucket while keeping image names in Images field.
func CreatePost(dynamodbClient dynamodbiface.DynamoDBAPI, tableName string, post Post) error {
	item, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = dynamodbClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// GetAllPostsOverview retrieves all posts from the DynamoDB table but only returns the ID, Title, Author, Summary, and Tag fields.
func GetAllPostsOverview(dynamodbClient dynamodbiface.DynamoDBAPI, tableName string) ([]Post, error) {
	input := &dynamodb.ScanInput{
		TableName:            aws.String(tableName),
		ProjectionExpression: aws.String("#id, #title, #author, #summary, #tags"),
		ExpressionAttributeNames: map[string]*string{
			"#id":      aws.String("id"),
			"#title":   aws.String("title"),
			"#author":  aws.String("author"),
			"#summary": aws.String("summary"),
			"#tags":    aws.String("tags"),
		},
	}

	result, err := dynamodbClient.Scan(input)
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0)
	for _, item := range result.Items {
		post := Post{}
		err = dynamodbattribute.UnmarshalMap(item, &post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
