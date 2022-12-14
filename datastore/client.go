package datastore

import (
	"context"

	"github.com/axatol/anywhere/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Client *s3.Client

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

	if _, err := Client.ListBuckets(ctx, &s3.ListBucketsInput{}); err != nil {
		panic(err)
	}

	config.Log.Infow("initialised datastore client",
		"bucket", config.Values.Datastore.Bucket,
		"prefix", config.Values.Datastore.PathPrefix,
	)
}
