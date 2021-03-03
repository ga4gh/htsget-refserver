package awsutils

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var mockCred = aws.Credentials{
	AccessKeyID:     "MOCKQ6HSRDFZ5JKZMOCK",
	SecretAccessKey: "M0CKSDr32mKyG/cbjceUj4IdiQEnlsKwNYtOT/V",
	SessionToken:    "MockFwoGZXIvYXdzEEcaDP9PvgOm1tNvvW2bFiLsARkAtzMOCKTbAX",
}

var tmpProfile = ""

func before() {
	val, exist := os.LookupEnv(AwsProfile)
	if exist {
		tmpProfile = val
	}
}

func after() {
	if len(tmpProfile) > 0 {
		_ = os.Setenv(AwsProfile, tmpProfile)
	}
}

// go test -run TestSetProfile ./internal/awsutils/ -v -count 1
func TestSetProfile(t *testing.T) {
	before()
	SetProfile("test")
	val, exist := os.LookupEnv(AwsProfile)
	assert.True(t, exist)
	assert.Equal(t, "test", val)
	after()
}

// go test -run TestSetRegion ./internal/awsutils/ -v -count 1
func TestSetRegion(t *testing.T) {
	before()
	SetRegion("us-east-2")
	val, exist := os.LookupEnv(AwsRegion)
	assert.True(t, exist)
	assert.Equal(t, "us-east-2", val)
	after()
}

// go test -run TestSetCredentials ./internal/awsutils/ -v -count 1
func TestSetCredentials(t *testing.T) {
	before()
	SetCredentials(mockCred)
	val, exist := os.LookupEnv(AwsSessionToken)
	assert.True(t, exist)
	assert.Equal(t, mockCred.SessionToken, val)
	after()
}

// go test -run TestUnsetCredentials ./internal/awsutils/ -v -count 1
func TestUnsetCredentials(t *testing.T) {
	before()
	SetCredentials(mockCred)
	assert.Equal(t, mockCred.AccessKeyID, os.Getenv(AwsAccessKeyId))
	UnsetCredentials()
	_, exist := os.LookupEnv(AwsAccessKeyId)
	assert.False(t, exist)
	after()
}

// go test -run TestGetCredentials ./internal/awsutils/ -v -count 1
func TestGetCredentials(t *testing.T) {
	before()
	SetProfile("")
	SetCredentials(mockCred)
	cred, err := GetCredentials()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("Using credentials source: " + cred.Source)
	assert.Equal(t, mockCred.AccessKeyID, cred.AccessKeyID)
	after()
}
