package queue

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MessageProcessorFunc func(ctx context.Context, messageID string, message Message) error

type Service struct {
	client *sqs.Client

	MessageProcessor MessageProcessorFunc

	QueueName string
	QueueUrl  string

	ChanMessage chan types.Message

	Mutex     sync.Mutex
	WaitGroup sync.WaitGroup
}

type QueueService interface {
	UpdateQueueUrl(ctx context.Context) error
	ConsumeMessages(ctx context.Context)
}

type EventType string

type Message struct {
	EventType EventType   `json:"event_type"`
	Data      interface{} `json:"data"`
}

type TopicNotification struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}
