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
import java.util.Objects;
import java.util.UUID;

public class Confluence {

    private static final String SCENARIO_SPACE_KEY_NAME = "confluence-space-key";
    private static final String SCENARIO_SPACE_NAME = "Space";
    private static final String SCENARIO_SPACE_HOMEPAGE_ID_KEY_NAME = "confluence-space-homepage-id-key";
    private static final String DRY_RUN_MODE = "dry-run-mode";

    public static String getScenarioSpaceKey() {
        return Objects.toString(ScenarioDataStore.get(SCENARIO_SPACE_KEY_NAME), "");
    }

    public static String getScenarioSpaceHomepageID() {
        return (String) ScenarioDataStore.get(SCENARIO_SPACE_HOMEPAGE_ID_KEY_NAME);
    }

    public static boolean isDryRun() {
        return (boolean) ScenarioDataStore.get(DRY_RUN_MODE);
    }

    @BeforeScenario
    public void setDryRunModeOff() {
        ScenarioDataStore.put(DRY_RUN_MODE, false);
    }

    @BeforeScenario
    public void setSpaceKeyName() {
        ScenarioDataStore.put(SCENARIO_SPACE_KEY_NAME, generateUniqueSpaceKeyName());
    }

    @BeforeScenario(tags = {"create-space-manually"})
    public void beforeScenario() {
        String spaceHomepageID = ConfluenceClient.createSpace(getScenarioSpaceKey(), SCENARIO_SPACE_NAME);
        ScenarioDataStore.put(SCENARIO_SPACE_HOMEPAGE_ID_KEY_NAME, spaceHomepageID);
    }

    @AfterScenario(tags = {"create-space-manually"})
    public void afterScenario() {
        ConfluenceClient.deleteSpace(getScenarioSpaceKey());
    }

    @Step("Activate dry run mode")
    public void activateDryRunMode() {
        ScenarioDataStore.put(DRY_RUN_MODE, true);
    }

    @Step("Space does not exist")
    public void assertSpaceDoesNotExist() {
        assertThat(ConfluenceClient.doesSpaceExist(getScenarioSpaceKey())).isFalse();
    }

    @Step("Published pages are: <table>")
    public void assertPublishedPages(Table expectedPages) throws Exception {
        int expectedTotal = expectedPages.getTableRows().size();
        assertConsoleSuccessOutput(expectedTotal - 1); // don't count the homepage as we don't publish it
        Space space = new Space(getScenarioSpaceKey());
        assertThat(space.totalPages()).isEqualTo(expectedTotal);
        for (TableRow row : expectedPages.getTableRows()) {
            String actualParentPageTitle = space.getParentPageTitle(row.getCell("title"));
            assertThat(actualParentPageTitle).isEqualTo(row.getCell("parent"));
        }
    }

    @Step("Specs <did|did not> get published")
    public void didPublishingOccur(String didPublishingOccur) throws IOException {
        boolean publishingOccurred = (didPublishingOccur.equalsIgnoreCase("did"));
        Space space = new Space(getScenarioSpaceKey());
        if (publishingOccurred) {
            new Console().outputContains("Success: published");
            assertThat(space.totalPages()).isGreaterThan(1);
        } else {
            new Console().outputDoesNotContain("Success: published");
            assertThat(space.totalPages()).isEqualTo(1);
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

    @Step("Manually delete the Confluence space homepage")
    public void manuallyDeleteTheConfluenceSpaceHomepage() {
       ConfluenceClient.deletePage(getScenarioSpaceHomepageID());
    }

    private void waitForNextMinuteToStart() throws InterruptedException {
        TimeUnit.SECONDS.sleep(60 - LocalTime.now().getSecond());
    }

    private String generateUniqueSpaceKeyName() {
        return UUID.randomUUID().toString().replace("-", "");
    }

}
