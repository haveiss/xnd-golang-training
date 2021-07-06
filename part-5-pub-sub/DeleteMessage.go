/*
   Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
// snippet-start:[sqs.go.delete_message]
package main

// snippet-start:[sqs.go.delete_message.imports]
import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// snippet-end:[sqs.go.delete_message.imports]

// GetQueueURL gets the URL of an Amazon SQS queue
// Inputs:
//     sess is the current session, which provides configuration for the SDK's service clients
//     queueName is the name of the queue
// Output:
//     If success, the URL of the queue and nil
//     Otherwise, an empty string and an error from the call to
func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMessages gets the messages from an Amazon SQS queue
// Inputs:
//     sess is the current session, which provides configuration for the SDK's service clients
//     queueURL is the URL of the queue
//     timeout is how long, in seconds, the message is unavailable to other consumers
// Output:
//     If success, the latest message and nil
//     Otherwise, nil and an error from the call to ReceiveMessage
func GetMessages(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	// snippet-start:[sqs.go.receive_messages.call]
	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   timeout,
	})
	// snippet-end:[sqs.go.receive_messages.call]
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

// DeleteMessage deletes a message from an Amazon SQS queue
// Inputs:
//     sess is the current session, which provides configuration for the SDK's service clients
//     queueURL is the URL of the queue
//     messageID is the ID of the message
// Output:
//     If success, nil
//     Otherwise, an error from the call to DeleteMessage
func DeleteMessage(sess *session.Session, queueURL *string, messageHandle *string) error {
	// snippet-start:[sqs.go.delete_message.call]
	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	})
	// snippet-end:[sqs.go.delete_message.call]
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// snippet-start:[sqs.go.delete_message.args]
	queue := os.Getenv("SQS_QUEUE_NAME")
	timeout := flag.Int64("t", 5, "How long, in seconds, that the message is hidden from others")

	if queue == "" {
		panic("You must supply the name of a queue (-q QUEUE)")
		return
	}
	// snippet-end:[sqs.go.delete_message.args]

	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	// snippet-start:[sqs.go.delete_message.sess]
	region := os.Getenv("SQS_AWS_REGION")
	endpoint := os.Getenv("SQS_ENDPOINT")

	sess, err := session.NewSession(&aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	})

	if err != nil {
		panic(err)
	}
	// snippet-end:[sqs.go.delete_message.sess]

	// Get URL of queue
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	msgResult, err := GetMessages(sess, queueURL, timeout)
	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return
	}

	for _, message := range msgResult.Messages {

		err = DeleteMessage(sess, queueURL, message.ReceiptHandle)
		if err != nil {
			fmt.Println("Got an error deleting the message:")
			fmt.Println(err)
			return
		}

		fmt.Println("[DELETED] Message ID:     " + *message.MessageId)
	}
}

// snippet-end:[sqs.go.delete_message]
