package com.example.controllers;

import com.example.MessageUtils;
import com.example.models.SimulateEventRequest;
import com.example.models.ValidatorGroups;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Profile;
import org.springframework.http.ResponseEntity;
import org.springframework.util.MimeType;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping
@RequiredArgsConstructor
@Profile("BOTH")
public class ControllerBoth {

    private static final MimeType MIME_AVRO = new MimeType("application", "*+avro");
    private static final MimeType MIME_JSON = new MimeType("application", "json");

    private final MessageUtils messageUtils;

    @PostMapping("/json")
    public ResponseEntity<?> postJSON(@Validated(value = ValidatorGroups.JsonValidator.class) @RequestBody SimulateEventRequest body) {
        body.getHeaders().put("content-type", MIME_JSON.toString());
        return new ResponseEntity<>(messageUtils.sendMessage(body.getTopic(), "json", messageUtils.createMessage(body.getPayload().getBytes(), body.getHeaders()), MIME_JSON));
    }

    @PostMapping("/avro")
    public ResponseEntity<?> postAVRO(@Validated(value = ValidatorGroups.AvroValidator.class) @RequestBody SimulateEventRequest body) {
        body.getHeaders().put("content-type", MIME_AVRO.toString());
        return new ResponseEntity<>(messageUtils.sendMessage(body.getTopic(), "avro", messageUtils.constructMessage(body), MIME_AVRO));
    }

}
