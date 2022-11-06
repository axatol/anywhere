package datastore

import (
	"context"

	"github.com/axatol/anywhere/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Client *s3.Client
var logger = config.Log.Named("datastore")

func Init(ctx context.Context) {
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               config.Values.Datastore.Host,
			SigningRegion:     "us-east-1",
			HostnameImmutable: true,
		}, nil
	})

	creds := credentials.NewStaticCredentialsProvider(
		config.Values.Datastore.User,
		config.Values.Datastore.Pass,
		"",
	)

	cfg := aws.Config{
		Region:                      "us-east-1",
		Credentials:                 creds,
		EndpointResolverWithOptions: endpointResolver,
	}

	Client = s3.NewFromConfig(cfg)
}
