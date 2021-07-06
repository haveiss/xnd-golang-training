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
// snippet-start:[sqs.go.send_message]
package main

// snippet-start:[sqs.go.send_message.imports]
import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// snippet-end:[sqs.go.send_message.imports]

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

// SendMsg sends a message to an Amazon SQS queue
// Inputs:
//     sess is the current session, which provides configuration for the SDK's service clients
//     queueURL is the URL of the queue
// Output:
//     If success, nil
//     Otherwise, an error from the call to SendMessage
func SendMsg(sess *session.Session, queueURL *string) error {
	// Create an SQS service client
	svc := sqs.New(sess)

	type exampleMessage struct {
		Name string
	}

	message := exampleMessage{"John Doe"}
	messageBody, _ := json.Marshal(message)

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
		QueueUrl:    queueURL,
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	// snippet-start:[sqs.go.send_message.args]
	queue := os.Getenv("SQS_QUEUE_NAME")

	if queue == "" {
		panic("You must supply the name of a queue (-q QUEUE)")
		return
	}
	// snippet-end:[sqs.go.send_message.args]

	// Create a session that gets credential values from environment variable
	// snippet-start:[sqs.go.send_message.sess]
	region := os.Getenv("SQS_AWS_REGION")
	endpoint := os.Getenv("SQS_ENDPOINT")

	sess, err := session.NewSession(&aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	})

	if err != nil {
		panic(err)
	}
	// snippet-end:[sqs.go.send_message.sess]

	// Get URL of queue
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	err = SendMsg(sess, queueURL)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}

	fmt.Println("Sent message to queue ")
}

// snippet-end:[sqs.go.send_message]
