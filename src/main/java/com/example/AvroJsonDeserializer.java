package com.example;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.confluent.kafka.serializers.KafkaAvroDeserializer;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.common.errors.UnsupportedForMessageFormatException;
import org.apache.kafka.common.header.Header;
import org.apache.kafka.common.header.Headers;
import org.apache.kafka.common.serialization.Deserializer;
import org.springframework.http.MediaType;
import org.springframework.kafka.support.serializer.DeserializationException;
import org.springframework.kafka.support.serializer.JsonDeserializer;
import org.springframework.util.MimeType;

import java.util.List;
import java.util.Map;

@Slf4j
public class AvroJsonDeserializer implements Deserializer<Object> {

    /**
     * List of potential values for the header key specifying the MediaType of the content
     */
    private static final List<String> CONTENT_TYPE_HEADER_VALUES = List.of("contentType", "content-type","Content-Type");

    private static final MimeType APPLICATION_AVRO = MimeType.valueOf("application/*+avro");
    private final KafkaAvroDeserializer kafkaAvroDeserializer;
    private final JsonDeserializer<Object> jsonDeserializer;
    /**
     * Default constructor initializes the {@link KafkaAvroDeserializer} instance which is then configured
     * on call to configure method
     */
    public AvroJsonDeserializer() {
        kafkaAvroDeserializer = new KafkaAvroDeserializer();
        jsonDeserializer = new JsonDeserializer<>(new ObjectMapper());
    }

    /**
     *
     * @param topic topic associated with the data
     * @param headers headers associated with the record; may be empty.
     * @param data serialized bytes; may be null; implementations are recommended to handle null by returning a value or null rather than throwing an exception.
     * @return data parsed as an object
     * @throws DeserializationException exception is caught with the org.springframework.kafka.support.serializer.ErrorHandlingDeserializer
     *  which is the method that actually invokes this method at runtime
     */
    @Override
    public Object deserialize(String topic, Headers headers, byte[] data) {
        MimeType contentType = this.extractContentTypeHeader(headers);
        if(MediaType.APPLICATION_JSON.isCompatibleWith(contentType)) {
            log.debug("JSON DESERIALIZER TRIGGERED");
            return jsonDeserializer.deserialize(topic, headers, data);
//            try {
//                return new ObjectMapper().readValue(data, Object.class);
//            } catch (IOException e) {
//                throw new DeserializationException("Failed to parse message as JSON", data, false, e);
//            }
        } else if (APPLICATION_AVRO.isCompatibleWith(contentType)) {
            log.debug("AVRO DESERIALIZER TRIGGERED");
            return kafkaAvroDeserializer.deserialize(topic, data);
        }
        throw new UnsupportedForMessageFormatException("Content type not supported: " + contentType);
    }

    /**
     * Used to close out any deserializers on exit
     */
    @Override
    public void close() {
        Deserializer.super.close();
        kafkaAvroDeserializer.close();
        jsonDeserializer.close();
    }

    /**
     * Used to configure prepared deserializers, e.g. {@link KafkaAvroDeserializer}
     * @param configs configs in key/value pairs
     * @param isKey whether is for key or value
     */
    @Override
    public void configure(Map<String, ?> configs, boolean isKey) {
        kafkaAvroDeserializer.configure(configs, isKey);
        Map<String, ?> mappy = Map.of("spring.json.trusted.packages", "*");
        jsonDeserializer.configure(mappy, isKey);
    }

    /**
     * Here to satisfy the interface requirements.
     *
     * @param s topic associated with the data
     * @param bytes serialized bytes; may be null; implementations are recommended to handle null by returning a value or null rather than throwing an exception.
     * @return
     */
    @Override
    public Object deserialize(String s, byte[] bytes) {
        // intentionally left blank
        log.error("Method should never be invoked");
        return null;
    }

    /**
     * Helper method to extract and parse the content type header
     * @param headers headers associated with the record; may be empty.
     * @return content type header converted to MimeType
     * @throws MissingEventHeaderException if content type header not found
     */
    private MimeType extractContentTypeHeader(Headers headers) {
        Header contentType;
        for(String headerValue: CONTENT_TYPE_HEADER_VALUES) {
            contentType = headers.lastHeader(headerValue);
            if (contentType != null) {
                return MimeType.valueOf(new String(contentType.value()).replace("\"", ""));
            }
        }
        throw new MissingEventHeaderException("Content Type header missing");
    }

}