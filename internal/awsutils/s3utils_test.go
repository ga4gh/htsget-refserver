package awsutils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type S3MockClient struct{}

func (client *S3MockClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	panic("implement me")
}

func (client *S3MockClient) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{
		ContentLength: int64(1111),
	}, nil
}

// go test -run TestHeadS3Object ./internal/awsutils/ -v -count 1
func TestHeadS3Object(t *testing.T) {
	contentLength, _ := HeadS3Object(S3Dto{
		ObjPath: "s3://does/not/matter.bam",
		Client:  &S3MockClient{},
	})
	fmt.Println("Content Length: ", contentLength)
	assert.Equal(t, int64(1111), contentLength)
}

// go test -run TestIntegrationHeadS3Object ./internal/awsutils/ -v -count 1
func TestIntegrationHeadS3Object(t *testing.T) {

	if os.Getenv(AwsProfile) != TestAwsProfileForIT {
		t.Skipf("[Skip] Required to setup and `export AWS_PROFILE=%s` in integration testing CI environment", TestAwsProfileForIT)
	}

	SetProfile(TestAwsProfileForIT)
	SetRegion("us-east-1")

	publicObjPath := "s3://giab/data/NA12878/Garvan_NA12878_HG001_HiSeq_Exome/project.NIST_NIST7035_H7AP8ADXX_TAAGGCGA_1_NA12878.bwa.markDuplicates.bam"
	expContentLength := int64(3020450026)

	contentLength, err := HeadS3Object(S3Dto{
		ObjPath: publicObjPath,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("Content Length: ", contentLength)
	assert.Equal(t, expContentLength, contentLength)
}
