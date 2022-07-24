package marshaller

import (
	"errors"
	"example.com/kafka-serializer-publisher/model"
)



const (
	Json model.EventContentType = "application/json"
	Avro model.EventContentType = "application/*+avro"
)

type Option func(handler Marshaller)

// Marshaller - Provides an interface to be able to register additional marshaller for different codecs.
type Marshaller interface {
	// RegisterMarshaller - Ability to register additional marshallers after the fact
	RegisterMarshaller(contentType model.EventContentType, marshaller Marshal)

	// Marshal - serializes the Publishable instance into the specified model.EventContentType using the appropriate
	// standards (e.g. Json, Avro)
	Marshal(p model.Publishable, eventContentType model.EventContentType) ([]byte, error)
}

// Marshal - generic definition of Marshal functions for our defined PUB/SUB model
// using Publishable and Event models
type Marshal interface {
	// Marshal - serializes the model.Publishable instance into the specified model.EventContentType using the appropriate
	// standards (e.g. Json, Avro)
	Marshal(p model.Publishable) ([]byte, error)

	// GetContentType - Get the content type of the marshaller (e.g. Json, Avro)
	GetContentType() model.EventContentType
}

// marshaller - Implements events.Marshal
type marshaller struct {
	avroExtension Marshal
	jsonExtension Marshal

	serializers map[model.EventContentType]Marshal
}

// New - returns an instance of Marshal which is configured to handle data of type:
// 		- Avro
// 		- Json
func New(marsh ...Marshal) Marshaller {
	marshes := &marshaller{
		serializers: make(map[model.EventContentType]Marshal),
	}
	// Append Marshaller to config
	for _, mar := range marsh {
		marshes.serializers[mar.GetContentType()] = mar
	}
	return marshes
}

func (d *marshaller) RegisterMarshaller(contentType model.EventContentType, marshaller Marshal) {
	d.serializers[contentType] = marshaller
}

// Marshal - serializes model.Publishable based on model.EventContentType
func (d *marshaller) Marshal(p model.Publishable, eventContentType model.EventContentType) ([]byte, error) {
	serializer := d.serializers[eventContentType]
	if serializer != nil {
		return serializer.Marshal(p)
	}
	return nil, errors.New("unsupported message content-type, value was: " + string(eventContentType))
}
