package projections

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/flypay/go-kit/v4/pkg/runtime"
)

func ResolveDynamoCompositeDb(table string, session *session.Session) (CompositeReadWriter, error) {
	return ResolveDynamoCompositeDbWithRuntimeConfig(table, session, runtime.DefaultConfig())
}

func ResolveDynamoCompositeDbWithRuntimeConfig(table string, session *session.Session, cfg runtime.Config) (CompositeReadWriter, error) {
	if cfg.ProjectionsInMemory {
		return &DynamoCompositeDB{}, nil
	}

	if cfg.ProjectionsLocalCreate {
		if err := createLocalDynamoCompositeTable(table, session); err != nil {
			return nil, err
		}
	}

	return NewDynamoCompositeDB(table, session)
}

type DynamoCompositeDB struct {
	dynamoDBCommon
}

func NewDynamoCompositeDB(table string, sess *session.Session) (*DynamoCompositeDB, error) {
	dynamoCommon, err := newDynamoDBCommon(table, sess)
	if err != nil {
		return nil, err
	}
	return &DynamoCompositeDB{
		dynamoDBCommon: dynamoCommon,
	}, nil
}

func (d DynamoCompositeDB) ReadComposite(ctx context.Context, id string, rangeKey *dynamodb.Condition, projections interface{}) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "ReadComposite", err) }()

	// Check that "projections" is a slice or a pointer to a slice
	check := projections
	if reflect.TypeOf(check).Kind() == reflect.Ptr {
		check = reflect.ValueOf(check).Elem().Interface()
	}
	if reflect.TypeOf(check).Kind() != reflect.Slice {
		return ErrInvalidProjectionsDestination
	}

	kc := map[string]*dynamodb.Condition{
		identifierField: {
			ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
			AttributeValueList: []*dynamodb.AttributeValue{
				{
					S: aws.String(id),
				},
			},
		},
	}

	if rangeKey != nil && len(rangeKey.AttributeValueList) > 0 {
		kc[rangeField] = rangeKey
	}

	queryInput := dynamodb.QueryInput{
		TableName:     aws.String(d.tableName),
		KeyConditions: kc,
	}

	if expr, ok := projectionExpression(projections); ok {
		queryInput.ProjectionExpression = expr.Projection()
		queryInput.ExpressionAttributeNames = expr.Names()
	}

	resultItems := []map[string]*dynamodb.AttributeValue{}
	for {
		results, err := d.driver.QueryWithContext(ctx, &queryInput)
		if err != nil {
			return err
		}
		resultItems = append(resultItems, results.Items...)

		if results.LastEvaluatedKey == nil {
			break
		}
		queryInput.ExclusiveStartKey = results.LastEvaluatedKey
	}
	if len(resultItems) < 1 {
		return ErrNotFound
	}
	return dynamodbattribute.UnmarshalListOfMaps(resultItems, &projections)
}

func (d DynamoCompositeDB) PatchComposite(ctx context.Context, id, rangeVal, path string, value interface{}, opts ...Option) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "PatchComposite", err) }()

	var cfg options
	cfg.apply(opts...)

	expr, err := buildPatchExpression(map[string]interface{}{path: value}, cfg)
	if err != nil {
		return err
	}

	return d.performPatchComposite(ctx, id, rangeVal, expr)
}

func (d DynamoCompositeDB) PatchFieldsComposite(ctx context.Context, id, rangeVal string, keyValues map[string]interface{}, opts ...Option) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "PatchFieldsComposite", err) }()

	var cfg options
	cfg.apply(opts...)

	expr, err := buildPatchExpression(keyValues, cfg)
	if err != nil {
		return err
	}

	return d.performPatchComposite(ctx, id, rangeVal, expr)
}

func (d DynamoCompositeDB) performPatchComposite(ctx context.Context, id, rangeVal string, expr expression.Expression) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			identifierField: {
				S: aws.String(id),
			},
			rangeField: {
				S: aws.String(rangeVal),
			},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              aws.String("NONE"),
	}

	_, err := d.driver.UpdateItemWithContext(ctx, input)
	var aerr awserr.Error
	if errors.As(err, &aerr) {
		if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return ErrNotFound
		}
	}

	return err
}

// DeleteComposite implements deleting only a single item
func (d DynamoCompositeDB) DeleteComposite(ctx context.Context, id, rangeVal string) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "DeleteComposite", err) }()
	_, err = d.driver.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: &d.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			identifierField: {S: aws.String(id)},
			rangeField:      {S: aws.String(rangeVal)},
		},
	})
	return err
}

func createLocalDynamoCompositeTable(table string, session *session.Session) error {
	dynamoClient := dynamodb.New(session)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(identifierField),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(rangeField),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(identifierField),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			}, {
				AttributeName: aws.String(rangeField),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		TableName: aws.String(table),
	}
	return createLocalTable(dynamoClient, input)
}
