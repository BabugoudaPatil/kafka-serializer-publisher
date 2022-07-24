package com.example.models;

public interface ValidatorGroups {
    interface AllValidator extends AvroValidator, JsonValidator {}
    interface AvroValidator {}
    interface JsonValidator {}
}
