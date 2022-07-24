package com.example;

public class MissingEventHeaderException extends RuntimeException {
    public MissingEventHeaderException(String msg) {
        super(msg);
    }
}
