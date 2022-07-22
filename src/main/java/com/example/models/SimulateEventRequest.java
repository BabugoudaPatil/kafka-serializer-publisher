package com.example.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
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
@JsonInclude(JsonInclude.Include.NON_NULL)
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
    @JsonSetter(nulls = Nulls.SET)
    @NotNull(message = "payload json is required. It can't be null or blank.", groups = ValidatorGroups.AvroValidator.class)
    private String payload;

    /**
     * A non-mandatory map of headers.
     */
    private Map<@NotBlank(message = "Header key name is required. It can't be empty", groups = {ValidatorGroups.JsonValidator.class, ValidatorGroups.AvroValidator.class}) String,
            @NotBlank(message = "Header value is required. It can't be empty", groups = {ValidatorGroups.JsonValidator.class, ValidatorGroups.AvroValidator.class}) String> headers;

    @JsonSetter(value = "payload")
    void setJsonPayload(JsonNode data) {
        this.payload = data.toString();
    }

}
