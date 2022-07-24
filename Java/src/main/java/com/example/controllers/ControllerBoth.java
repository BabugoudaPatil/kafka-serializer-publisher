package com.example.controllers;

import com.example.MessageUtils;
import com.example.models.SimulateEventRequest;
import com.example.models.ValidatorGroups;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Profile;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import static com.example.MessageUtils.MIME_AVRO;
import static com.example.MessageUtils.MIME_JSON;

@RestController
@RequestMapping
@RequiredArgsConstructor
@Profile("default")
@Validated
public class ControllerBoth {

    private final MessageUtils messageUtils;

    @PostMapping("/json")
    public ResponseEntity<?> postJSON(@Validated(value = ValidatorGroups.JsonValidator.class) @RequestBody SimulateEventRequest body) {
        byte[] payload = {};
        if (body.getPayload() != null) {
            payload = body.getPayload().getBytes();
        }
        return new ResponseEntity<>(messageUtils.sendMessage(body.getTopic(), "json", messageUtils.createMessage(payload, body.getHeaders(), MIME_JSON), MIME_JSON));
    }

    @PostMapping("/avro")
    public ResponseEntity<?> postAVRO(@Validated(value = ValidatorGroups.AvroValidator.class) @RequestBody SimulateEventRequest body) {
        return new ResponseEntity<>(messageUtils.sendMessage(body.getTopic(), "avro", messageUtils.constructAvroMessage(body), MIME_AVRO));
    }

}
