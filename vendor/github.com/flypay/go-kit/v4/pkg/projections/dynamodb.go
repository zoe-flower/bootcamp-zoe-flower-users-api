package projections

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/flypay/go-kit/v4/pkg/runtime"
)

const (
	identifierField = "Identifier"
	rangeField      = "Range"
	expirationField = "ExpiresAt"
	defaultRegion   = "eu-west-1"
)

type dynamoDBCommon struct {
	metrics   instrumenter
	driver    dynamodbiface.DynamoDBAPI
	tableName string
}

func newDynamoDBCommon(table string, sess *session.Session) (dynamoDBCommon, error) {
	var zero dynamoDBCommon
	if sess == nil {
		var err error
		sess, err = MakeEnvAwareSession("")
		if err != nil {
			return zero, err
		}
	}

	client := dynamodb.New(sess)
	xray.AWS(client.Client)
	dbMetrics := registerDynamoMetrics()

	return dynamoDBCommon{
		tableName: table,
		driver:    client,
		metrics:   dbMetrics,
	}, nil
}

// Write implements the `Writer` interface backed by dynamodb.
// Under the hood we use `PutItem` to overwrite any previously stored state.
func (d dynamoDBCommon) Write(ctx context.Context, data interface{}, expires *time.Time) (err error) {
	now := time.Now()
	defer func() { d.metrics.MetricsReq(now, "Write", err) }()

	var item map[string]*dynamodb.AttributeValue

	item, err = dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	if _, ok := item[identifierField]; !ok {
		return ErrIdentifierFieldMissing
	}

	if expires != nil {
		timestamp := strconv.Itoa(int(expires.Unix()))

		item[expirationField] = &dynamodb.AttributeValue{
			N: &timestamp,
		}
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(d.tableName),
	}

	if field, v := sequenceField(data); field != nil {
		conditionalValue := v.FieldByName(field.Name).Int()

		cond := expression.Name(field.Name).
			LessThanEqual(expression.Value(conditionalValue)).
			Or(expression.AttributeNotExists(expression.Name(field.Name)))
		builder := expression.NewBuilder().WithCondition(cond)
		expr, _ := builder.Build()

		input.ConditionExpression = expr.Condition()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	_, err = d.driver.PutItemWithContext(ctx, input)

	var e *dynamodb.ConditionalCheckFailedException
	if errors.As(err, &e) {
		// May want to improve this to wrap the conditional check failed exception as
		// it may contain useful information which is currently being lost.
		return ErrDataIntegrityMismatch
	}

	return err
}

// TableAndDriver provides access to the underlying table name and dynamo driver.
// This is to allow running non-standard queries against Dynamo while still
// having all the resources created and initialised properly.
//
// In order to access this, we can do 'projections.ReadWriter.(*DynamoDB).TableAndDriver()'
func (d dynamoDBCommon) TableAndDriver() (string, dynamodbiface.DynamoDBAPI) {
	return d.tableName, d.driver
}

type DynamoDB struct {
	dynamoDBCommon
}

func ResolveDynamoDb(table string, session *session.Session) (ReadWriter, error) {
	return ResolveDynamoDbWithRuntimeConfig(table, session, runtime.DefaultConfig())
}

func ResolveDynamoDbWithRuntimeConfig(table string, session *session.Session, cfg runtime.Config) (ReadWriter, error) {
	if cfg.ProjectionsInMemory {
		return &MemDB{}, nil
	}

	if cfg.ProjectionsLocalCreate {
		if err := createLocalDynamoTable(table, session); err != nil {
			return nil, err
		}
	}

	return NewDynamoDB(table, session)
}

//counterfeiter:generate -o internal/mocks/fake_dynamodbapi.go github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface.DynamoDBAPI

// NewDynamoDB returns a new instance of the dynamodb abstraction.
// It uses the provided session or if nil creates a new session in the default region.
// The client will be wrapped with XRay tracing which means any contexts you provide to the abstraction API must contain an XRay segment.
func NewDynamoDB(table string, sess *session.Session) (*DynamoDB, error) {
	dynamoCommon, err := newDynamoDBCommon(table, sess)
	if err != nil {
		return nil, err
	}
	return &DynamoDB{
		dynamoDBCommon: dynamoCommon,
	}, nil
}

// Read implements the `Reader` interface backed by dynamodb.
// We use `GetItem` to fetch the whole object.
func (d DynamoDB) Read(ctx context.Context, id string, projection interface{}) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "Read", err) }()

	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			identifierField: {
				S: aws.String(id),
			},
		},
	}

	if expr, ok := projectionExpression(projection); ok {
		getInput.ProjectionExpression = expr.Projection()
		getInput.ExpressionAttributeNames = expr.Names()
	}

	var result *dynamodb.GetItemOutput

	result, err = d.driver.GetItemWithContext(ctx, getInput)
	if err != nil {
		return err
	}

	if result.Item == nil {
		return ErrNotFound
	}

	return dynamodbattribute.UnmarshalMap(result.Item, &projection)
}

func (d DynamoDB) Patch(ctx context.Context, id, path string, value interface{}, opts ...Option) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "Patch", err) }()

	var cfg options
	cfg.apply(opts...)

	expr, err := buildPatchExpression(map[string]interface{}{path: value}, cfg)
	if err != nil {
		return err
	}

	return d.performPatch(ctx, id, expr)
}

func (d DynamoDB) PatchFields(ctx context.Context, id string, keyValues map[string]interface{}, opts ...Option) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "PatchFields", err) }()

	var cfg options
	cfg.apply(opts...)

	expr, err := buildPatchExpression(keyValues, cfg)
	if err != nil {
		return err
	}

	return d.performPatch(ctx, id, expr)
}

func (d DynamoDB) performPatch(ctx context.Context, id string, expr expression.Expression) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			identifierField: {
				S: aws.String(id),
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

func (d DynamoDB) Delete(ctx context.Context, id string) (err error) {
	start := time.Now()
	defer func() { d.metrics.MetricsReq(start, "Delete", err) }()

	_, err = d.driver.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			identifierField: {
				S: aws.String(id),
			},
		},
	})
	return err
}

func MakeEnvAwareSession(name string) (*session.Session, error) {
	return MakeEnvAwareSessionWithRuntimeConfig(name, runtime.DefaultConfig())
}

func MakeEnvAwareSessionWithRuntimeConfig(name string, cfg runtime.Config) (*session.Session, error) {
	if cfg.ProjectionsEndpoint != "" {
		return newCustomEndpointSession(cfg.ProjectionsEndpoint)
	}

	return newDefaultSession()
}

func newDefaultSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(defaultRegion),
	})
}

func newCustomEndpointSession(endpoint string) (*session.Session, error) {
	return session.NewSession(
		aws.NewConfig().
			WithCredentials(credentials.NewStaticCredentials("flyt", "flyt1234", "")).
			WithEndpoint(endpoint).
			WithRegion("eu-west-1"),
	)
}

type instrumenter interface {
	MetricsReq(startTime time.Time, reqType string, err error)
}

func createLocalDynamoTable(table string, session *session.Session) error {
	dynamoClient := dynamodb.New(session)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(identifierField),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(identifierField),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		TableName: aws.String(table),
	}
	return createLocalTable(dynamoClient, input)
}

func createLocalTable(dynamoClient *dynamodb.DynamoDB, input *dynamodb.CreateTableInput) error {
	input.SetBillingMode(dynamodb.BillingModePayPerRequest)
	_, err := dynamoClient.CreateTable(input)
	if aerr, ok := err.(awserr.Error); ok { // nolint:errorlint
		switch aerr.Code() {
		case dynamodb.ErrCodeTableAlreadyExistsException, dynamodb.ErrCodeResourceInUseException:
			return nil
		default:
			return aerr
		}
	}
	return err
}

// sequenceField returns the struct field that should be used for sequencing if defined
func sequenceField(data interface{}) (*reflect.StructField, reflect.Value) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	var field *reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		if f := t.Field(i); f.Tag.Get("sequencing") == "true" {
			field = &f
			break
		}
	}

	return field, v
}

func buildPatchExpression(keyValues map[string]interface{}, cfg options) (expression.Expression, error) {
	var (
		updateBuilder expression.UpdateBuilder
		conditions    []expression.ConditionBuilder
	)

	for key, value := range keyValues {
		updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(value))

		if !cfg.createOnNotExist {
			conditions = append(conditions, expression.AttributeExists(expression.Name(key)))
		}
	}

	builder := expression.NewBuilder().WithUpdate(updateBuilder)

	if len(conditions) > 0 {
		builder = builder.WithCondition(requireAllConditions(conditions))
	}

	return builder.Build()
}

func requireAllConditions(updateConditions []expression.ConditionBuilder) expression.ConditionBuilder {
	condition := updateConditions[0]

	if len(updateConditions) > 1 {
		condition = condition.And(updateConditions[1], updateConditions[2:]...)
	}

	return condition
}

func projectionExpression(projection interface{}) (expression.Expression, bool) {
	// Pull out the correct field name by using the dynamodbattribute.Marshal
	// method on the incoming projection.
	check := projection

	// If the projection is a pointer, get the concrete type
	if check != nil && reflect.TypeOf(check).Kind() == reflect.Ptr {
		check = reflect.ValueOf(check).Elem().Interface()
	}

	// If the projection is a slice, get the concrete underlying type
	if check != nil && reflect.TypeOf(check).Kind() == reflect.Slice {
		check = reflect.New(reflect.TypeOf(check).Elem()).Elem().Interface()
	}

	// If the projection is still a pointer, then zero it out and get the
	// concrete type
	if check != nil && reflect.TypeOf(check).Kind() == reflect.Ptr {
		check = reflect.Zero(reflect.TypeOf(check).Elem()).Interface()
	}

	if check == nil {
		return expression.Expression{}, false
	}

	attrs, _ := dynamodbattribute.Marshal(check)
	if attrs.M == nil {
		return expression.Expression{}, false
	}

	// Sort them so we can consistently create name expressions
	var attributes []string

	for k := range attrs.M {
		attributes = append(attributes, k)
	}

	sort.Strings(attributes)

	// Create the name expressions, this will help us leverage DynamoDB to
	// create the field aliases for us.
	var nameExpressions []expression.NameBuilder

	for _, attribute := range attributes {
		nameExpressions = append(nameExpressions, expression.Name(attribute))
	}

	proj := expression.NamesList(nameExpressions[0], nameExpressions[1:]...)

	expr, _ := expression.NewBuilder().
		WithProjection(proj).
		Build()

	return expr, true
}
