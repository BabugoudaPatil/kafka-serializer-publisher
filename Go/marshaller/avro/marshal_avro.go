package avro

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"example.com/kafka-serializer-publisher/go-kafka-avro"
	"example.com/kafka-serializer-publisher/marshaller"
	"example.com/kafka-serializer-publisher/model"
)

const (
	contentType = "application/*+avro"
)

type extension struct {
	schemaRegistryClient *kafka.CachedSchemaRegistryClient
	config               SchemaRegistry
}

// New - returns an instance of marshaller.Marshal which is configured to handle data of type marshaller.Avro
func New(cfg SchemaRegistry) marshaller.Marshal {
	if len(cfg.Hosts) == 0 {
		cfg.Hosts = []string{"http://localhost:8081"}
	}
	return &extension{schemaRegistryClient: kafka.NewCachedSchemaRegistryClient(cfg.Hosts), config: cfg}
}

// GetContentType - Returns content type for marshaller
func (ext *extension) GetContentType() model.EventContentType {
	return contentType
}

// Marshal - used for publishing events where payload conforms to marshaller.Avro content-type
func (ext *extension) Marshal(p model.Publishable) ([]byte, error) {
	if p.Subject() == nil {
		return nil, errors.New("subject is nil")
	}
	id, codec, err := ext.schemaRegistryClient.GetLatestSchema(*p.Subject())
	if err != nil {
		return nil, err
	}

	value, _ := json.Marshal(p.Payload())
	// Validates that obj does not violate the codec
	native, _, err := codec.NativeFromTextual(value)
	if err != nil {
		return nil, err
	}
	// Convert native Go form to binary Avro data
	binaryValue, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return nil, err
	}

	binaryMsg := &encoder{
		SchemaID: id,
		Content:  binaryValue,
	}
	return binaryMsg.encode()
}

// transcode used to convert one struct to another, specifically a map[string]interface{} to from struct
func (ext *extension) transcode(in, out interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(in); err != nil {
		return err
	}
	if err := json.NewDecoder(buf).Decode(out); err != nil {
		return err
	}
	return nil
}

// encoder encodes schemaId as magic bytes into Avro message
type encoder struct {
	SchemaID int
	Content  []byte
}

// Notice: the Confluent schema registry has special requirements for the Avro serialization rules,
// not only need to serialize the specific content, but also attach the Schema ID and Magic Byte.
// Ref: https://docs.confluent.io/current/schema-registry/serializer-formatter.html#wire-format
func (enc *encoder) encode() ([]byte, error) {
	var binaryMsg []byte
	// Confluent serialization format version number; currently always 0.
	binaryMsg = append(binaryMsg, byte(0))
	// 4-byte schema ID as returned by Schema Registry
	binarySchemaId := make([]byte, 4)
	binary.BigEndian.PutUint32(binarySchemaId, uint32(enc.SchemaID))
	binaryMsg = append(binaryMsg, binarySchemaId...)
	// Avro serialized data in Avro's binary encoding
	binaryMsg = append(binaryMsg, enc.Content...)
	return binaryMsg, nil
}

// Length of schemaId and Content.
func (enc *encoder) Length() int {
	return 5 + len(enc.Content)
}
