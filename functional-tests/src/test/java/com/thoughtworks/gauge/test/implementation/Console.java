package com.thoughtworks.gauge.test.implementation;

import static com.thoughtworks.gauge.test.common.GaugeProject.getCurrentProject;
import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;

import com.thoughtworks.gauge.Step;

public class Console {

    @Step({ "Output contains <message>", "The error message <message> should be output" })
    public void outputContains(String message) throws IOException {
        String output = getCurrentProject().getStdOut();
        assertThat(output).contains(message);
    }

    @Step({ "Output contains <message> <message2> <message3>" })
    public void outputContainsMessages(String message, String message2, String message3) throws IOException {
        String output = getCurrentProject().getStdOut();
        assertThat(output).contains(message + message2 + message3);
    }

    @Step({ "Output is <message>" })
    public void outputIs(String message) throws IOException {
        String output = getCurrentProject().getStdOut();
        assertThat(output.trim()).isEqualTo(message);
    }
}
