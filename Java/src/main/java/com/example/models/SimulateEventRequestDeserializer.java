package com.example.models;

import com.fasterxml.jackson.databind.util.StdConverter;

import java.util.HashMap;

public class SimulateEventRequestDeserializer extends StdConverter<SimulateEventRequest, SimulateEventRequest> {
    @Override
    public SimulateEventRequest convert(SimulateEventRequest simulateEventRequest) {
        if (simulateEventRequest.getHeaders() == null) {
            simulateEventRequest.setHeaders(new HashMap<>());
        }
        return simulateEventRequest;
    }
}
