package ddbmoose

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DdbMooseCreate(awsRegion string) (*DdbMoose, error) {
	if len(awsRegion) == 0 {
		return nil, errors.New("AWS region not sent")
	}

	awsConfig, failed := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if failed != nil {
		return nil, errors.New("Error creating upload settings from AWS. Details: " + failed.Error())
	}

	clientDynamodb := dynamodb.NewFromConfig(awsConfig)

	if _, failed := clientDynamodb.ListTables(context.TODO(), &dynamodb.ListTablesInput{}); failed != nil {
		fmt.Println("error")
		return nil, errors.New("AWS Region sent invalid. Details: " + failed.Error())
	} else {
		return &DdbMoose{clientDynamodb: clientDynamodb}, nil
	}
}

func (ddb *DdbMoose) SetTable(tableName string) error {
	if len(tableName) == 0 {
		return errors.New("table name is empty")
	}

	if _, failed := ddb.clientDynamodb.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}); failed != nil {
		return errors.New("Table not found. Details: " + failed.Error())
	} else {
		ddb.tableName = tableName
		return nil
	}
}

func (ddb *DdbMoose) Save(data interface{}) (interface{}, error) {
	mapData, failed := attributevalue.MarshalMap(data)

	if failed != nil {
		return nil, errors.New("Format sent to invalid database. Details: " + failed.Error())
	}

	created, failed := ddb.clientDynamodb.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(ddb.tableName),
		Item:      mapData,
	})

	if failed != nil {
		return nil, errors.New("Failed to save data. Details: " + failed.Error())
	} else if failed = attributevalue.UnmarshalMap(created.Attributes, &data); failed != nil {
		return nil, errors.New(failed.Error())
	} else {
		return data, nil
	}
}

func (ddb *DdbMoose) FindWithFilters(filters *[]DdbMooseFilter) ([]interface{}, error) {
	expr, failed := generateExpression(filters)

	if failed != nil {
		return nil, errors.New("Failed to generate the filters in the database. Details: " + failed.Error())
	}

	result, failed := ddb.clientDynamodb.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(ddb.tableName),
		Select:                    types.SelectAllAttributes,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if failed != nil {
		return nil, errors.New("Failed fetching data. Details: " + failed.Error())
	}

	var data []interface{}
	if failed = attributevalue.UnmarshalListOfMaps(result.Items, &data); failed != nil {
		return nil, errors.New("Failed to convert database data to view. Details: " + failed.Error())
	}

	return data, nil
}

func (ddb *DdbMoose) Delete(fieldKeyPrimary string, valueKeyPrimary types.AttributeValue) error {
	if result, failed := ddb.clientDynamodb.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &ddb.tableName,
		Key: map[string]types.AttributeValue{
			fieldKeyPrimary: valueKeyPrimary,
		},
	}); failed != nil {
		return errors.New("Failed fetching data. Details: " + failed.Error())
	} else if result.Item == nil {
		return errors.New("data not found to delete")
	} else if _, failed = ddb.clientDynamodb.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &ddb.tableName,
		Key: map[string]types.AttributeValue{
			fieldKeyPrimary: valueKeyPrimary,
		},
	}); failed != nil {
		return errors.New("Failed to delete data. Details: " + failed.Error())
	} else {
		return nil
	}
}
