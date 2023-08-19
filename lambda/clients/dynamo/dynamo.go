package dynamo

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UserStorageDB interface {
	GetUser(username string) (*User, error)
	ValidateRefreshToken(username, refreshToken string) bool
	AddUserToDB(username, password, hashedToken string) error
	UpdateUserToken(username, token string) error
}

const (
	TABLE_NAME = "user-table-name"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DynamoDBClient struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBClient() UserStorageDB {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	return &DynamoDBClient{
		db: db,
	}
}

func (db *DynamoDBClient) GetUser(username string) (*User, error) {
	result, err := db.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	var user User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DynamoDBClient) ValidateRefreshToken(username string, refreshToken string) bool {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}

	result, err := db.db.GetItem(input)
	if err != nil {
		// TODO: Handle error
		return false
	}

	if len(result.Item) == 0 {
		// TODO: Handle error
		return false
	}

	storedHashedToken := aws.StringValue(result.Item["hashedToken"].S)
	return storedHashedToken == refreshToken
}

func (db *DynamoDBClient) AddUserToDB(username string, password string, hashedToken string) error {
	item := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
			"password": {
				S: aws.String(password),
			},
			"hashedToken": {
				S: aws.String(hashedToken),
			},
		},
		TableName: aws.String(TABLE_NAME),
	}

	_, err := db.db.PutItem(item)
	return err

}

func (db *DynamoDBClient) UpdateUserToken(username, token string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":token": {
				S: aws.String(token),
			},
		},
		UpdateExpression: aws.String("SET token = :token"),
	}

	_, err := db.db.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
