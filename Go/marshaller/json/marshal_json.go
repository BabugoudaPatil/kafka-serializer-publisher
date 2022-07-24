package json

import (
	"encoding/json"
	"example.com/kafka-serializer-publisher/marshaller"
	"example.com/kafka-serializer-publisher/model"
)

const (
	contentType = "application/json"
)

// extension - Implements marshaller.Marshal
type extension struct {
}

// New - returns an instance of marshaller.Marshal which is configured to handle data of type marshaller.Json
func New() marshaller.Marshal {
	return &extension{}
}

// GetContentType - Returns content type
func (d *extension) GetContentType() model.EventContentType {
	return contentType
}

// Marshal - used for publishing events where payload conforms to marshaller.Json content-type
func (d *extension) Marshal(p model.Publishable) ([]byte, error) {
	_ = json.RawMessage{}
	return json.Marshal(p.Payload())
}
