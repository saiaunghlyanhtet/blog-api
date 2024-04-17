package post

import (
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type Post struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Summary     string   `json:"summary"`
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
	CreatedDate string   `json:"createdDate"`
}

// CreatePost creates a new post in the DynamoDB table.
func CreatePost(dynamodbClient dynamodbiface.DynamoDBAPI, tableName string, req events.APIGatewayProxyRequest) error {
	log.Println("QueryParameters ", req.QueryStringParameters)
	var post Post
	id := uuid.New().String()
	post.ID = id
	post.Author = req.QueryStringParameters["author"]
	post.Title = req.QueryStringParameters["title"]
	post.Summary = req.QueryStringParameters["summary"]
	post.Content = req.QueryStringParameters["content"]
	post.Images = []string{req.QueryStringParameters["images"]}
	post.Tags = []string{req.QueryStringParameters["tags"]}
	post.CreatedDate = req.QueryStringParameters["createdDate"]

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

// GetPostById retrieves a post from the DynamoDB table by ID.
var (
	s3Session  *session.Session
	bucketName string
)

func InitializeS3SessionAndBucket(session *session.Session, name string) {
	s3Session = session
	bucketName = name
}

func GetPostById(dynamodbClient dynamodbiface.DynamoDBAPI, tableName, id string, createdDate string) (*Post, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"createdDate": {
				S: aws.String(createdDate),
			},
		},
	}

	result, err := dynamodbClient.GetItem(input)
	if err != nil {
		return nil, err
	}

	post := Post{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		return nil, err
	}

	for i, imageName := range post.Images {
		url, err := generatePresignedURL(imageName)
		if err != nil {
			return nil, err
		}
		post.Images[i] = url
	}

	return &post, nil
}

func generatePresignedURL(key string) (string, error) {
	s3Client := s3.New(s3Session)
	expiration := 15 * time.Minute

	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", err
	}

	return url, nil
}

// DeletePost deletes a post from the DynamoDB table by ID.
func DeletePost(dynamodbClient dynamodbiface.DynamoDBAPI, tableName, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := dynamodbClient.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
