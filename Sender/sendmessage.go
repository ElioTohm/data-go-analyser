package Sender

import (
	"context"
	"data-go-analyser/data"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
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

func TransformMessageToSQSMessage(messages []*data.Customer) (sqsMessage []types.SendMessageBatchRequestEntry) {
	for _, message := range messages {
		messageID := message.CustomerReference + time.Now().Format("20060102150405")
		messageJson, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		sqsMessage = append(sqsMessage, types.SendMessageBatchRequestEntry{
			Id: &messageID,
			MessageAttributes: map[string]types.MessageAttributeValue{
				"customer_reference": {
					DataType:    aws.String("String"),
					StringValue: &message.CustomerReference,
				},
			},
			MessageBody: aws.String(string(messageJson)),
		})
	}
	return sqsMessage
}
