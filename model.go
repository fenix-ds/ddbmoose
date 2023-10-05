package ddbmoose

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DdbMoose struct {
	tableName      string
	clientDynamodb *dynamodb.Client
}

type DdbMooseFilter struct {
	Field           string
	Operation       TypeFilter
	Value           interface{}
	LogicalOperator TypeLogicalOperator
}
