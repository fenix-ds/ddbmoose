package ddbmoose

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func generateExpression(filters *[]DdbMooseFilter) (*expression.Expression, error) {
	var filtersScan expression.ConditionBuilder

	for _, filter := range *filters {
		switch filter.LogicalOperator {
		case TloAnd:
			filtersScan = filtersScan.And(generateFilterScan(filter))
		case TloOr:
			filtersScan = filtersScan.Or(generateFilterScan(filter))
		case TloNone:
			filtersScan = generateFilterScan(filter)
		default:
			return nil, errors.New("invalid logical operator type")
		}
	}

	result, failed := expression.NewBuilder().WithFilter(filtersScan).Build()

	return &result, failed
}

func generateFilterScan(filter DdbMooseFilter) expression.ConditionBuilder {
	switch filter.Operation {
	case TfrEqual:
		return expression.Name(filter.Field).Equal(expression.Value(filter.Value))
	case TfrNotEqual:
		return expression.Name(filter.Field).NotEqual(expression.Value(filter.Value))
	case TfrContains:
		return expression.Name(filter.Field).Contains(fmt.Sprint(filter.Value))
	case TfrBeginsWith:
		return expression.Name(filter.Field).BeginsWith(fmt.Sprint(filter.Value))
	case TfrGreaterThan:
		return expression.Name(filter.Field).GreaterThan(expression.Value(filter.Value))
	case TfrGreaterThanEqual:
		return expression.Name(filter.Field).GreaterThanEqual(expression.Value(filter.Value))
	case TfrLessThan:
		return expression.Name(filter.Field).LessThan(expression.Value(filter.Value))
	case TfrLessThanEqual:
		return expression.Name(filter.Field).LessThanEqual(expression.Value(filter.Value))
	case TfrBetween:
		list := strings.Split(fmt.Sprint(filter.Value), ",")
		listInt := listStringToListInt(list)
		return expression.Name(filter.Field).Between(expression.Value(listInt[0]), expression.Value(listInt[1]))
	default:
		return expression.ConditionBuilder{}
	}
}

func listStringToListInt(list []string) [2]int64 {
	var listInt [2]int64

	for count, item := range list {
		if count < 2 {
			inte, failed := strconv.ParseInt(item, 6, 12)

			if failed != nil {
				listInt[count] = 0
			} else {
				listInt[count] = inte
			}
		}
	}

	return listInt
}
