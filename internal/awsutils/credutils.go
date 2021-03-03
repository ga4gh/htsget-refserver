package awsutils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
	"strconv"
)

func SetProfile(profileName string) {
	_ = os.Setenv(AwsProfile, profileName)
}

func SetRegion(region string) {
	_ = os.Setenv(AwsRegion, region)
}

func SetCredentials(cred aws.Credentials) {
	_ = os.Setenv(AwsAccessKeyId, cred.AccessKeyID)
	_ = os.Setenv(AwsSecretAccessKey, cred.SecretAccessKey)
	_ = os.Setenv(AwsSessionToken, cred.SessionToken)
	_ = os.Setenv(AwsSessionTokenExpiration, strconv.FormatInt(cred.Expires.Unix(), 10))
}

func UnsetCredentials() {
	_ = os.Unsetenv(AwsAccessKeyId)
	_ = os.Unsetenv(AwsSecretAccessKey)
	_ = os.Unsetenv(AwsSessionToken)
}

func GetCredentials() (*aws.Credentials, error) {
	cfg, cfgErr := config.LoadDefaultConfig(context.TODO())
	if cfgErr != nil {
		return nil, cfgErr
	}

	cred, credErr := cfg.Credentials.Retrieve(context.TODO())

	return &cred, credErr
}
