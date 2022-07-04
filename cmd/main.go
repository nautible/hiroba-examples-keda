package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	sqsSvc    *sqs.SQS
	accountId string
	region    string
	sqsUrl    string
)

func receiveMessages(chn chan<- *sqs.Message) {
	for {
		output, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqsUrl),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(10),
		})
		if err != nil {
			panic(err)
		}
		for _, message := range output.Messages {
			chn <- message
		}
	}
}

func main() {
	accountId = os.Getenv("ACCOUNT_ID")
	region = os.Getenv("REGION")
	sqsUrl = fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/demo-keda-queue", region, accountId)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		panic(err)
	}

	sqsSvc = sqs.New(sess)
	messages := make(chan *sqs.Message, 1)
	go receiveMessages(messages)

	for message := range messages {
		fmt.Println("receive message ...")
		fmt.Println(*message.Body)
		sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sqsUrl),
			ReceiptHandle: message.ReceiptHandle,
		})
	}
}
