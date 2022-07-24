package model

import "context"

type EventContentType string

type Publisher interface {
	Disconnect()
	Send(ctx context.Context, contentType EventContentType, publishable Publishable) error
}

type Publishable interface {
	ID() string
	Headers() map[string]string
	Payload() interface{}
	// Subject - nullable, specific to AVRO based events.
	//
	// The subject in Schema Registry used to look up the avro specification
	Subject() *string
	// Topic - non-null
	//
	// The topic to publish events to
	Topic() string
}
