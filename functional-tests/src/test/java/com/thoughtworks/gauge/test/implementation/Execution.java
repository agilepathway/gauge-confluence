package com.thoughtworks.gauge.test.implementation;

import static com.thoughtworks.gauge.test.common.GaugeProject.getCurrentProject;

import com.thoughtworks.gauge.Step;
import com.thoughtworks.gauge.test.common.ExecutionSummary;
import com.thoughtworks.gauge.test.common.ExecutionSummaryAssert;

public class Execution {

    private String getFormattedProcessOutput() {
        return "\n*************** Process output start************\n" + getCurrentProject().getLastProcessStdout()
                + "\n*************** Process output end************\n";
    }

    private ExecutionSummaryAssert assertOn(ExecutionSummary summary, boolean result) {
        return ExecutionSummaryAssert.assertThat(summary).withFailMessage(getFormattedProcessOutput())
                .hasSuccess(result);
    }

    @Step("Publish Confluence Documentation for the current project")
    public void publishConfluenceDocumentationForCurrentProject() throws Exception {
        assertOn(getCurrentProject().publishConfluenceDocumentation(), true);
    }

    @Step("Publish Confluence Documentation for the current project and assert <did|did not> succeed")
    public void publishConfluenceDocumentationForCurrentProjectAndAssertResult(String didOrDidNotSucceed) throws Exception {
        boolean result = (didOrDidNotSucceed.equalsIgnoreCase("did"));
        assertOn(getCurrentProject().publishConfluenceDocumentation(), result);
    }

    @Step("Publish Confluence Documentation with no <variable> configured and assert <did|did not> succeed")
    public void publishConfluenceDocumentationWithConfigVarUnset(String configVar, String didOrDidNotSucceed) throws Exception {
        new ProjectInit().projectInit();
        boolean result = (didOrDidNotSucceed.equalsIgnoreCase("did"));
        assertOn(getCurrentProject().publishConfluenceDocumentationWithConfigVarUnset(configVar), result);
    }

    @Step("Publish Confluence Documentation for two projects")
    public void publishConfluenceDocumentationForTwoProjects() throws Exception {
        assertOn(getCurrentProject().publishConfluenceDocumentationForTwoProjects(), true);
    }
}
