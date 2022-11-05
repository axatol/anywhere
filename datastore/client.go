package datastore

import (
	"context"
	"fmt"

	"github.com/tunes-anywhere/anywhere/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Client *s3.Client
var logger = config.Log.Named("datastore")

func Init(ctx context.Context) {
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               config.Config.Datastore.Host,
			SigningRegion:     "us-east-1",
			HostnameImmutable: true,
		}, nil
	})

	fmt.Println(
		"User", config.Config.Datastore.User,
		"Pass", config.Config.Datastore.Pass,
	)

	cfg := aws.Config{
		Region: "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider(
			config.Config.Datastore.User,
			config.Config.Datastore.Pass,
			"session",
		),
		EndpointResolver: staticResolver,
	}

	Client = s3.NewFromConfig(cfg)
}
