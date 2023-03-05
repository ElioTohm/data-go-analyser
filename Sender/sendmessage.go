package Sender

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSSendMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)

	SendMessageBatch(ctx context.Context,
		params *sqs.SendMessageBatchInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error)
}

func GetQueueURLInterface(c context.Context, api SQSSendMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func SendMsgInterface(c context.Context, api SQSSendMessageAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}

func SendMsgBulkInterface(c context.Context, api SQSSendMessageAPI, input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	return api.SendMessageBatch(c, input)
}

func SendSuccessMessage(sqsClient *sqs.Client, queueName *string, messages []types.SendMessageBatchRequestEntry) {

	result, err := GetQueueURLInterface(context.TODO(), sqsClient, &sqs.GetQueueUrlInput{
		QueueName: queueName,
	})
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	sMinput := &sqs.SendMessageBatchInput{
		Entries:  messages,
		QueueUrl: queueURL,
	}

	_, err = SendMsgBulkInterface(context.TODO(), sqsClient, sMinput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}

}
