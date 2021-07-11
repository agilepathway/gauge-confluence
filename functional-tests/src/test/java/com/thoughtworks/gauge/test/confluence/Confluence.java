package com.thoughtworks.gauge.test.confluence;

import com.thoughtworks.gauge.BeforeScenario;
import com.thoughtworks.gauge.Step;
import com.thoughtworks.gauge.Table;
import com.thoughtworks.gauge.TableRow;
import com.thoughtworks.gauge.AfterScenario;
import com.thoughtworks.gauge.datastore.ScenarioDataStore;
import com.thoughtworks.gauge.test.implementation.Console;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.time.LocalTime;
import java.util.concurrent.TimeUnit;
import java.util.UUID;

public class Confluence {

    private static final String SCENARIO_SPACE_KEY_NAME = "confluence-space-key";
    private static final String SCENARIO_SPACE_NAME = "Space";

    public static String getScenarioSpaceKey() {
        return (String) ScenarioDataStore.get(SCENARIO_SPACE_KEY_NAME);
    }

    @BeforeScenario
    public void BeforeScenario() {
        ScenarioDataStore.put(SCENARIO_SPACE_KEY_NAME, generateUniqueSpaceKeyName());
        ConfluenceClient.createSpace(getScenarioSpaceKey(), SCENARIO_SPACE_NAME);
    }

    @AfterScenario
    public void AfterScenario() {
        ConfluenceClient.deleteSpace(getScenarioSpaceKey());
    }

    @Step("Published pages are: <table>")
    public void createSpec(Table expectedPages) throws Exception {
        int expectedTotal = expectedPages.getTableRows().size();
        assertConsoleSuccessOutput(expectedTotal - 1); // don't count the homepage as we don't publish it
        Space space = new Space(getScenarioSpaceKey());
        assertThat(space.totalPages()).isEqualTo(expectedTotal);
        for (TableRow row : expectedPages.getTableRows()) {
            String actualParentPageTitle = space.getParentPageTitle(row.getCell("title"));
            assertThat(actualParentPageTitle).isEqualTo(row.getCell("parent"));
        }
    }

    private void assertConsoleSuccessOutput(int totalPages) throws IOException {
        new Console().outputContains(
                String.format("Success: published %d specs and directory pages to Confluence", totalPages));
    }

    @Step("Manually add a page to the Confluence space")
    public void manuallyAddPageToConfluenceSpace() throws InterruptedException {
        // the page needs to be added at a later minute than when the last publish ran
        waitForNextMinuteToStart();
        ConfluenceClient.createPage(getScenarioSpaceKey());
        TimeUnit.SECONDS.sleep(2); // give Confluence time to process the added page
    }

    private void waitForNextMinuteToStart() throws InterruptedException {
        TimeUnit.SECONDS.sleep(60 - LocalTime.now().getSecond());
    }

    private String generateUniqueSpaceKeyName() {
        return UUID.randomUUID().toString().replace("-", "");
    }

}
