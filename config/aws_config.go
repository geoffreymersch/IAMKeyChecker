package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func SetupIAMClient(profileName string) (iamClient *iam.Client) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profileName), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), 5)
	}))
	if err != nil {
		log.Println(err)
	}

	iamClient = iam.NewFromConfig(cfg)

	return iamClient
}
