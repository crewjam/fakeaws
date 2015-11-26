package fakedynamodb

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FakeDynamoDBTest struct {
}

var _ = Suite(&FakeDynamoDBTest{})

func (s *FakeDynamoDBTest) TestStuff(c *C) {
	fakeDB, err := New()
	c.Assert(err, IsNil)
	defer fakeDB.Close()
	c.Assert(fakeDB.Port, Not(Equals), 0)

	db := dynamodb.New(session.New(), fakeDB.Config)

	_, err = db.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String("frob"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("Key"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	c.Assert(err, IsNil)
}
