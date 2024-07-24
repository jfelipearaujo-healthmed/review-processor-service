package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func NewService(queueName string, config aws.Config, messageProcessor MessageProcessorFunc) QueueService {
	return &Service{
		client: sqs.NewFromConfig(config),

		MessageProcessor: messageProcessor,
		QueueName:        queueName,

		ChanMessage: make(chan types.Message, 10),

		Mutex:     sync.Mutex{},
		WaitGroup: sync.WaitGroup{},
	}
}

func (svc *Service) UpdateQueueUrl(ctx context.Context) error {
	result, err := svc.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(svc.QueueName),
	})
	if err != nil {
		return err
	}

	svc.QueueUrl = *result.QueueUrl

	return nil
}

func (svc *Service) ConsumeMessages(ctx context.Context) {
	output, err := svc.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &svc.QueueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     30,
	})
	if err != nil {
		slog.ErrorContext(ctx, "error receiving message from queue", "queue_url", svc.QueueUrl, "error", err)
		return
	}

	svc.WaitGroup.Add(len(output.Messages))

	for _, message := range output.Messages {
		go svc.processMessage(ctx, message)
	}

	svc.WaitGroup.Wait()
}

func (svc *Service) processMessage(ctx context.Context, message types.Message) {
	defer svc.WaitGroup.Done()
	svc.Mutex.Lock()

	slog.InfoContext(ctx, "message received", "message_id", *message.MessageId)

	var notification TopicNotification

	if err := json.Unmarshal([]byte(*message.Body), &notification); err != nil {
		slog.ErrorContext(ctx, "error unmarshalling message", "error", err)
	} else {
		if notification.Type != "Notification" {
			slog.ErrorContext(ctx, "message is not a notification", "message_id", *message.MessageId)
		} else {
			var request Message

			if err := json.Unmarshal([]byte(notification.Message), &request); err != nil {
				slog.ErrorContext(ctx, "error unmarshalling message", "message_id", *message.MessageId, "error", err)
			} else {
				slog.InfoContext(ctx, "message unmarshalled", "request", request)
				if err := svc.MessageProcessor(ctx, notification.MessageId, request); err != nil {
					slog.ErrorContext(ctx, "error processing message", "message_id", *message.MessageId, "error", err)
				}
			}
		}
	}

	if err := svc.deleteMessage(ctx, message); err != nil {
		slog.ErrorContext(ctx, "error deleting message", "message_id", *message.MessageId, "error", err)
	}

	svc.Mutex.Unlock()
}

func (svc *Service) deleteMessage(ctx context.Context, message types.Message) error {
	_, err := svc.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &svc.QueueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return err
	}

	return nil
}
