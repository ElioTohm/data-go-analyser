package main

import (
	"context"
	"data-go-analyser/Reader"
	"data-go-analyser/Sender"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ztrue/tracerr"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, s3Event events.S3Event) (string, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		tracerr.PrintSourceColor(err)
		return "", err
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	sqsClient := sqs.NewFromConfig(sdkConfig)
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("Bucket = %s, Key = %s \n", s3.Bucket.Name, s3.Object.Key)
		// load s3 file in memory
		reader, err := Reader.DownloadFile(s3Client, s3.Bucket.Name, s3.Object.Key)
		if err != nil {
			tracerr.PrintSourceColor(err)
			return fmt.Sprintln("DownloadFile::", err), err
		}
		// transform svc lines to structs/objects
		constructedData, _ := Reader.ConstructData(Reader.ReadFile(reader))
		// analyse the struct and return useful data
		analysedData := Reader.AnalyseData(constructedData)

		Sender.SendSuccessMessage(sqsClient, aws.String("elio-s3-test"), Sender.TransformMessageToSQSMessage(analysedData))
	}
	return "200", nil
}
