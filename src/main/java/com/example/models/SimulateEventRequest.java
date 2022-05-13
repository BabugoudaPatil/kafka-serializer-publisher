package com.example.models;

import com.fasterxml.jackson.annotation.*;
import com.fasterxml.jackson.databind.JsonNode;
import lombok.*;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import java.util.Map;

@Builder
@AllArgsConstructor
@EqualsAndHashCode
@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
@Data
@NoArgsConstructor
public class SimulateEventRequest {
    /**
     * A mandatory topic name field.
     */
    @NotBlank(message = "topic field is required. It can't be null or blank.", groups = ValidatorGroups.AllValidator.class)
    private String topic;

    /**
     * A conditionally mandatory avroSource name field.
     */
    @NotBlank(message = "avroSource field is required. It can't be null or blank.", groups = ValidatorGroups.AvroValidator.class)
    private String avroSource;

    /**
     * A mandatory payload specified as a JSON object
     */
    @JsonIgnore
    @NotEmpty(message = "payload json is required. It can't be null or blank.", groups = ValidatorGroups.AllValidator.class)
    private String payload;

    /**
     * A mandatory map of headers.
     */
    @Singular
    @JsonProperty("headers")
    @NotEmpty(message = "Headers are required. It can't be empty")
    private Map<@NotBlank(message = "Header key name is required. It can't be empty") String,
            @NotBlank(message = "Header value is required. It can't be empty", groups = ValidatorGroups.AllValidator.class) String> headers;

    @JsonSetter("payload")
    void setJsonPayload(JsonNode data) {
        this.payload = data.toString();
    }
}
