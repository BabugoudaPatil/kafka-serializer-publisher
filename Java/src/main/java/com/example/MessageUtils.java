package com.example;

import com.example.models.SimulateEventRequest;
import converter.JsonAvroConverter;
import converter.ProcessingException;
import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient;
import io.confluent.kafka.schemaregistry.client.SchemaMetadata;
import io.confluent.kafka.schemaregistry.client.rest.RestService;
import io.confluent.kafka.schemaregistry.client.rest.exceptions.RestClientException;
import lombok.extern.slf4j.Slf4j;
import org.apache.avro.Schema;
import org.apache.avro.generic.GenericData;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.stream.function.StreamBridge;
import org.springframework.http.HttpStatus;
import org.springframework.integration.support.MessageBuilder;
import org.springframework.lang.Nullable;
import org.springframework.messaging.Message;
import org.springframework.stereotype.Service;
import org.springframework.util.MimeType;

import java.io.IOException;
import java.util.Map;

@Slf4j
@Service
public class MessageUtils {

    public static final MimeType MIME_JSON = new MimeType("application", "json");
    public static final MimeType MIME_AVRO = new MimeType("application", "*+avro");
    public static final String CONTENT_TYPE = "content-type";

    private final StreamBridge streamBridge;
    private final String schemaUrl;

    public MessageUtils(StreamBridge streamBridge,
                        @Value("${spring.cloud.stream.kafka.binder.producer-properties.schema.registry.url}")
                        String schemaUrl) {
        this.streamBridge = streamBridge;
        this.schemaUrl = schemaUrl;
    }

    public HttpStatus sendMessage(String topic, @Nullable String binding, Message<?> body, MimeType mimeType) {
        boolean sent = streamBridge.send(topic, binding, body, mimeType);
        if (sent) {
            return HttpStatus.ACCEPTED;
        }
        return HttpStatus.INTERNAL_SERVER_ERROR;
    }

    public <T> Message<T> createMessage(final T event, Map<String, String> headers, MimeType mimeType) {
        return MessageBuilder.withPayload(event)
                .copyHeaders(headers)
                .setHeaderIfAbsent(CONTENT_TYPE, mimeType.toString())
                .build();
    }

    public Message<GenericData.Record> constructAvroMessage(SimulateEventRequest request) {
        JsonAvroConverter avroConverter = new JsonAvroConverter();
        GenericData.Record record = avroConverter.convertToGenericDataRecord((request.getPayload()).getBytes(), getSchemaForType(request.getAvroSource()));
        return createMessage(record, request.getHeaders(), MIME_AVRO);
    }

    private Schema getSchemaForType(String type) {
        RestService restService = new RestService(schemaUrl);
        CachedSchemaRegistryClient cachedSchemaRegistryClient = new CachedSchemaRegistryClient(restService, 1);
        try {
            SchemaMetadata metadata = cachedSchemaRegistryClient.getLatestSchemaMetadata(type);
            return (Schema) cachedSchemaRegistryClient.getSchemaById(metadata.getId()).rawSchema();
        } catch (IOException | RestClientException e) {
            log.error("Exception while parsing schema", e);
            throw new ProcessingException(e.getMessage());
        }
    }

}