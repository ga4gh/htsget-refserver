package awsutils

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	log "github.com/sirupsen/logrus"
)

type S3ClientApi interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type S3Dto struct {
	ObjPath string
	Client  S3ClientApi
}

func (dto *S3Dto) getBucketAndKey() (string, string) {
	objPath := dto.ObjPath
	trimmedPath := strings.TrimPrefix(objPath, S3Proto)
	bucketName := strings.Split(trimmedPath, "/")[0]
	objKeyName := strings.TrimPrefix(trimmedPath, bucketName+"/")
	return bucketName, objKeyName
}

func (dto *S3Dto) NewS3Client() S3ClientApi {
	if dto.Client != nil {
		return dto.Client
	}

	defaultCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Debugf("error in LoadDefaultConfig: %v", err)
		return nil
	}
	return s3.NewFromConfig(defaultCfg)
}

func HeadS3Object(dto S3Dto) (int64, error) {
	client := dto.NewS3Client()
	bucketName, objKeyName := dto.getBucketAndKey()
	log.Debugf("client in s3Head %v", client)
	headResp, herr := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objKeyName),
	})
	log.Debugf("response from head %v", headResp)
	if herr != nil {
		log.Debugf("error in HeadObject: %v", herr)
		return 0, herr
	}
	return headResp.ContentLength, nil
}
