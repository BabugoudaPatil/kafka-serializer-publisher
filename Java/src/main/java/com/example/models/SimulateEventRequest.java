package com.example.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import java.util.Map;

@Builder
@AllArgsConstructor
@EqualsAndHashCode
@JsonIgnoreProperties(ignoreUnknown = true)
@Data
@NoArgsConstructor
@JsonDeserialize(converter = SimulateEventRequestDeserializer.class)
public class SimulateEventRequest {
    /**
     * A mandatory topic name field.
     */
    @NotBlank(message = "topic field is required. It can't be null or blank.", groups = {ValidatorGroups.JsonValidator.class, ValidatorGroups.AvroValidator.class})
    private String topic;

    /**
     * A conditionally mandatory avroSource name field.
     */
    @NotBlank(message = "avroSource field is required. It can't be null or blank.", groups = ValidatorGroups.AvroValidator.class)
    private String avroSource;

    /**
     * A payload specified as a JSON object
     */
    private JsonNode payload;

    /**
     * A non-mandatory map of headers.
     */
    private Map<@NotBlank(message = "Header key name is required. It can't be empty", groups = {ValidatorGroups.JsonValidator.class, ValidatorGroups.AvroValidator.class}) String,
            @NotBlank(message = "Header value is required. It can't be empty", groups = {ValidatorGroups.JsonValidator.class, ValidatorGroups.AvroValidator.class}) String> headers;

    @NotNull(message = "payload json is required. It can't be null or blank.", groups = ValidatorGroups.AvroValidator.class)
    public String getPayload() {
        System.out.println("Payload: " + this.payload);
        if (this.payload == null || this.payload.isNull()) {
            System.out.println("Payload Inner: " + this.payload);
            return null;
        }
        return this.payload.toString();
    }

}
