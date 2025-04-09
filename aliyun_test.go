package tests

import (
	"os"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/dm-20151123/client"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func TestAliyunOssUpload(t *testing.T) {

	const (
		OSS_BUCKET_REGION   = "cn-hongkong"
		OSS_BUCKET_ENDPOINT = "https://oss-cn-hongkong.aliyuncs.com"
		OSS_BUCKET_NAME     = "haiwell-hongkong"
		OSS_OBJECT_KEY      = "assets/a.txt"
		OSS_UPLOAD_FILEPATH = "./assets/a.txt"
	)

	accessKeyId := os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")

	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		t.Fatalf("Failed to create credentials provider: %v", err)
	}

	clientOptions := []oss.ClientOption{
		oss.SetCredentialsProvider(&provider),
		oss.Region(OSS_BUCKET_REGION),
		oss.AuthVersion(oss.AuthV4),
	}

	client, err := oss.New(OSS_BUCKET_ENDPOINT, accessKeyId, accessKeySecret, clientOptions...)
	if err != nil {
		t.Fatalf("Failed to create OSS client: %v", err)
	}

	bucket, err := client.Bucket(OSS_BUCKET_NAME)
	if err != nil {
		t.Fatalf("Failed to get bucket: %v", err)
	}

	err = bucket.PutObjectFromFile(OSS_OBJECT_KEY, OSS_UPLOAD_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to put object from file: %v", err)
	}

	t.Logf("%s\n", "File uploaded successfully.")

}

func TestAliyunEmailSend(t *testing.T) {
	const (
		EMAIL_ENDPOINT          = "dm.aliyuncs.com"
		EMAIL_ACCOUNT           = "noreply@mail.synwell.net"
		EMAIL_ADDRESS_TYPE      = 1
		EMAIL_REPLAY_TO_ADDRESS = false
		EMAIL_TO_ADDRESS        = "messy_things@outlook.com"
		EMAIL_SUBJECT           = "Test Email"
		EMAIL_HTML_BODY         = "<h1>Test Email</h1>"
	)
	accessKeyId := os.Getenv("ALIYUN_EMAIL_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_EMAIL_ACCESS_KEY_SECRET")

	endpoint := EMAIL_ENDPOINT

	c, err := client.NewClient(&openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
		Endpoint:        &endpoint,
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	accountName := EMAIL_ACCOUNT
	addressType := int32(EMAIL_ADDRESS_TYPE)
	replyToAddress := EMAIL_REPLAY_TO_ADDRESS
	toAddress := EMAIL_TO_ADDRESS
	subject := EMAIL_SUBJECT
	htmlBody := EMAIL_HTML_BODY

	singleSendMailRequest := &client.SingleSendMailRequest{
		AccountName:    &accountName,
		AddressType:    &addressType,
		ReplyToAddress: &replyToAddress,
		ToAddress:      &toAddress,
		Subject:        &subject,
		HtmlBody:       &htmlBody,
	}

	resp, err := c.SingleSendMail(singleSendMailRequest)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Logf("Email sent successfully: %v\n", resp)
}
