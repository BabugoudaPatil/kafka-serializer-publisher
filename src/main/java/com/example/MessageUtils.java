package com.example;

import com.example.models.SimulateEventRequest;
import converter.JsonAvroConverter;
import converter.ProcessingException;
import io.confluent.kafka.schemaregistry.ParsedSchema;
import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient;
import io.confluent.kafka.schemaregistry.client.rest.RestService;
import io.confluent.kafka.schemaregistry.client.rest.exceptions.RestClientException;
import lombok.RequiredArgsConstructor;
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
import java.util.List;
import java.util.Map;
import java.util.Optional;

@Slf4j
@Service
@RequiredArgsConstructor
public class MessageUtils {

    private final StreamBridge streamBridge;

    @Value("${spring.cloud.stream.kafka.binder.producer-properties.schema.registry.url}")
    String schemaUrl;

    public HttpStatus sendMessage(String topic, @Nullable String binding, Message<?> body, MimeType mimeType) {
        boolean sent = streamBridge.send(topic, binding, body, mimeType);
        if (sent) {
            return HttpStatus.ACCEPTED;
        }
        return HttpStatus.INTERNAL_SERVER_ERROR;
    }

    public <T> Message<T> createMessage(final T event, Map<String, String> headers) {
        return MessageBuilder.withPayload(event)
                .copyHeaders(headers)
                .build();
    }

    public <T extends GenericData.Record> Message<GenericData.Record> constructMessage(SimulateEventRequest request) {
        String avroClassName = request.getAvroSource();
        log.info("avroclassname:" + avroClassName);
        JsonAvroConverter avroConverter = new JsonAvroConverter();
        GenericData.Record record = avroConverter.convertToGenericDataRecord((request.getPayload()).getBytes(), getSchemaForType(avroClassName));
        return createMessage(record, request.getHeaders());
    }

    private Schema getSchemaForType(String type) {
        log.info("schemaUrl:{}", schemaUrl);
        RestService restService = new RestService(schemaUrl);
        CachedSchemaRegistryClient cachedSchemaRegistryClient = new CachedSchemaRegistryClient(restService, 1);
        Optional<ParsedSchema> schema = null;
        try {
            List<ParsedSchema> schemas = cachedSchemaRegistryClient.getSchemas("com", false, true);
            log.info("Parsed Schema");
            for (ParsedSchema schema1 : schemas) {
                log.debug(schema1.name() + ",");
            }
            schema = schemas.stream().filter(s -> type.equalsIgnoreCase(s.name())).findFirst();
            log.debug(String.valueOf(schemas));
        } catch (IOException | RestClientException e) {
            log.error("Exception while parsing schema", e);
            throw new ProcessingException(e.getMessage());
        }

        log.info("Schema found on schema registry.");
        log.info("schema:" + schema);
        return (Schema) schema.get().rawSchema();

    }

}