package Sender

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
)

type Message struct {
	Type              string
	CustomerReference string
	NumberOfOrders    int
	TotalAmountSpent  float32
}

type ErrorMessage struct {
	Type               string
	Customer_reference string
	Order_reference    string
	Message            string
}

func TransformMessageToSQSMessage(messages []*Message) (sqsMessage []types.SendMessageBatchRequestEntry) {
	for _, message := range messages {
		messageID := message.CustomerReference + time.Now().Format("20060102150405")
		sqsMessage = append(sqsMessage, types.SendMessageBatchRequestEntry{
			Id: &messageID,
			MessageAttributes: map[string]types.MessageAttributeValue{
				"type": {
					DataType:    aws.String("String"),
					StringValue: aws.String("customer_message"),
				},
				"customer_reference": {
					DataType:    aws.String("String"),
					StringValue: &message.CustomerReference,
				},
				"number_of_orders": {
					DataType:    aws.String("Number"),
					StringValue: aws.String(strconv.Itoa(message.NumberOfOrders)),
				},
				"total_amount_spent": {
					DataType:    aws.String("Number"),
					StringValue: aws.String(strconv.Itoa(int(message.TotalAmountSpent))),
				},
			},
			MessageBody: &message.CustomerReference,
		})
	}
	return sqsMessage
}
