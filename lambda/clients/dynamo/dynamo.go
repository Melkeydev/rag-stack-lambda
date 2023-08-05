package dynamo

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME = "user-table-name"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUserFromDynamoDB(username string) (*User, error) {
	sess := session.Must(session.NewSession())
	dynamoDB := dynamodb.New(sess)

	result, err := dynamoDB.GetItem(&dynamodb.GetItemInput{
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

func AddUserToDynamoDB(username string, password string, hashedToken string) error {
	sess := session.Must(session.NewSession())
	dynamoDB := dynamodb.New(sess)

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

	_, err := dynamoDB.PutItem(item)
	return err

}

func UpdateUserTokenInDynamoDB(username, token string) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

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

	_, err := svc.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
