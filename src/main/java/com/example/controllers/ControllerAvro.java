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

import static com.example.MessageUtils.MIME_AVRO;

@RestController
@RequestMapping
@RequiredArgsConstructor
@Profile("AVRO")
public class ControllerAvro {

    private final MessageUtils messageUtils;

    @PostMapping("/avro")
    public ResponseEntity<?> postAVRO(@Validated(value = ValidatorGroups.AvroValidator.class) @RequestBody SimulateEventRequest body) {
        return new ResponseEntity<>(messageUtils.sendMessage(body.getTopic(), null, messageUtils.constructAvroMessage(body), MIME_AVRO));
    }


}
