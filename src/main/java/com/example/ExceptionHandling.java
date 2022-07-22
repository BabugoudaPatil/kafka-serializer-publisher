package com.example;

import lombok.Builder;
import lombok.Data;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.lang.NonNull;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@ControllerAdvice
public class ExceptionHandling extends ResponseEntityExceptionHandler {

    private static final Logger LOG = LoggerFactory.getLogger(ExceptionHandling.class);

    @ExceptionHandler(Throwable.class)
    protected ResponseEntity<Object> handleGenericThrowable(Throwable ex, WebRequest request) {
        Map<String, Object> ret = new HashMap<>();
        ret.put("code", HttpStatus.INTERNAL_SERVER_ERROR.getReasonPhrase());
        ret.put("path", request.getDescription(false));
        ret.put("message", ex.getMessage());
        LOG.error("Error Handler", ex);
        return ResponseEntity.badRequest().body(ret);
    }

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
