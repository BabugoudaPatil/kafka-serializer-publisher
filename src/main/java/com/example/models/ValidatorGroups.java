package com.example.models;

public interface ValidatorGroups {
    public interface AllValidator extends AvroValidator, JsonValidator {}
    public interface AvroValidator {}
    public interface JsonValidator {}
}
