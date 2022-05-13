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
@Profile("AVRO")
public class ControllerAvro {

    private static final MimeType MIME_AVRO = new MimeType("application", "*+avro");

    private final MessageUtils messageUtils;

    @PostMapping("/avro")
    public ResponseEntity<?> postAVRO(@Validated(value = ValidatorGroups.AvroValidator.class) @RequestBody SimulateEventRequest body) {
        body.getHeaders().put("content-type", MIME_AVRO.toString());
        return new ResponseEntity<>(messageUtils.sendMessage(body.getTopic(), null, messageUtils.constructMessage(body), MIME_AVRO));
    }


}
