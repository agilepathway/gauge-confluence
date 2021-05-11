package com.thoughtworks.gauge.test.implementation;

import java.util.List;

import com.thoughtworks.gauge.Table;

public class Steps {

    public static Table example() {
        List<String> headers = List.of("step text", "implementation");
        Table steps = new Table(headers);
        steps.addRow(List.of("A step", "implementation"));
        return steps;
    }

}
