# Fake AWS

This package contains golang wrappers for AWS services that are useful to stub
out for local testing. Currently implemented are:

 - DynamoDB using the Amazon-provided DynamoDBLocal
 - S3 using the s3test package from goamz.
 