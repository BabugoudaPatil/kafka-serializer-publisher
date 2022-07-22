package com.example;

import lombok.Builder;
import lombok.Data;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.lang.NonNull;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.context.request.WebRequest;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@ControllerAdvice
public class ExceptionHandling extends ResponseEntityExceptionHandler {

    @Override
    @NonNull
    protected ResponseEntity<Object> handleMethodArgumentNotValid(MethodArgumentNotValidException ex,
                                                                  HttpHeaders headers, HttpStatus status,
                                                                  WebRequest request) {
        Map<String, Object> ret = new HashMap<>();
        ret.put("code", HttpStatus.BAD_REQUEST.getReasonPhrase());
        ret.put("path", request.getDescription(false));

        List<Error> errors = new ArrayList<>();
        ex.getBindingResult().getAllErrors().forEach((error) -> {
            errors.add(Error.builder()
                    .object(error.getObjectName())
                    .field(((FieldError) error).getField())
                    .rejectedValue(((FieldError) error).getRejectedValue())
                    .message(error.getDefaultMessage())
                    .build());
        });
        ret.put("errors", errors);
        return ResponseEntity.badRequest().body(ret);
    }


    @Data
    @Builder
    public static class Error {
        private String object;
        private String field;
        private Object rejectedValue;
        private String message;
    }

}
