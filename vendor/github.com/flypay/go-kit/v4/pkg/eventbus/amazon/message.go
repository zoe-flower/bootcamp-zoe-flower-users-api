package amazon

import (
	"time"
)

// SNSMessage represents an SNS message body that's been sent to an SQS Queue
// The full spec can be found here: https://docs.aws.amazon.com/sns/latest/dg/sns-message-and-json-formats.html
type SNSMessage struct {
	Type              string                      `json:"Type"`
	Token             string                      `json:"Token"`
	MessageID         string                      `json:"MessageId"`
	TopicArn          string                      `json:"TopicArn"`
	Subject           string                      `json:"Subject"`
	Message           string                      `json:"Message"`
	Timestamp         time.Time                   `json:"Timestamp"`
	SignatureVersion  string                      `json:"SignatureVersion"`
	Signature         string                      `json:"Signature"`
	SigningCertURL    string                      `json:"SigningCertURL"`
	UnsubscribeURL    string                      `json:"UnsubscribeURL"`
	SubscribeURL      string                      `json:"SubscribeURL"`
	MessageAttributes map[string]MessageAttribute `json:"MessageAttributes"`
}

type MessageAttribute struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

// ConfirmSubscriptionMessage is a subset of a full SNS Message which is handed to a SubscriptionConfirmer handler.
type ConfirmSubscriptionMessage struct {
	TopicARN string
	Token    string
}
