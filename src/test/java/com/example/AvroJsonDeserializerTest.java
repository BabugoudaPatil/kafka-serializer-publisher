package com.example;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.http.MediaType;
import org.springframework.util.MimeType;

import java.util.UUID;

class AvroJsonDeserializerTest {

    @Test
    void testParseMimeType() {
        Assertions.assertEquals(MediaType.APPLICATION_JSON, MimeType.valueOf("application/json"));
    }
    @Test
    void testParseMimeTypeException() {
        Assertions.assertThrows(IllegalArgumentException.class, () -> new MimeType("application/json"));
    }
    @Test
    void testParseMimeTypeEThing() {
        Assertions.assertTrue(MimeType.valueOf("application/avro").isCompatibleWith(MimeType.valueOf("application/*+avro")));
        Assertions.assertTrue(MimeType.valueOf("application/*+avro").isCompatibleWith(MimeType.valueOf("application/avro")));
    }

}
