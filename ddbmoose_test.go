package ddbmoose

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type Person struct {
	Id   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

func TestDdbmoose(t *testing.T) {
	ddbMoose, failed := DdbMooseCreate("sa-east-1")

	if failed != nil {
		t.Error(failed.Error())
	}

	failed = ddbMoose.SetTable("ddbmoose")

	if failed != nil {
		t.Error(failed.Error())
	}

	id := uuid.New().String()
	resultData, failed := ddbMoose.Save(Person{
		Id:   id,
		Name: "Mauricio",
	})

	if failed != nil {
		t.Error(failed.Error())
	}

	fmt.Println(resultData)

	resultData, failed = ddbMoose.Save(Person{
		Id:   id,
		Name: "Mauricio Martins",
	})

	if failed != nil {
		t.Error(failed.Error())
	}

	fmt.Println(resultData)

	resultFilters, failed := ddbMoose.FindWithFilters(&[]DdbMooseFilter{
		{Field: "name", Operation: TfrContains, Value: "Qau", LogicalOperator: TloNone},
	})

	if failed != nil {
		t.Error(failed.Error())
	}

	fmt.Println(resultFilters)

	resultFilters, failed = ddbMoose.FindWithFilters(&[]DdbMooseFilter{
		{Field: "name", Operation: TfrContains, Value: "Mau", LogicalOperator: TloNone},
	})

	if failed != nil {
		t.Error(failed.Error())
	}

	fmt.Println(resultFilters)

	failed = ddbMoose.Delete("id", &types.AttributeValueMemberS{Value: id})

	if failed != nil {
		t.Error(failed.Error())
	}
}
