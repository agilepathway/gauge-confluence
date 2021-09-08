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
    private static final String CONFLUENCE_USERNAME = "confluence-username";
    private static final String CONFLUENCE_TOKEN = "confluence-token";

    public static String getScenarioSpaceKey() {
        return Objects.toString(ScenarioDataStore.get(SCENARIO_SPACE_KEY_NAME), "");
    }

    public static String getScenarioSpaceHomepageID() {
        return (String) ScenarioDataStore.get(SCENARIO_SPACE_HOMEPAGE_ID_KEY_NAME);
    }

    public static boolean isDryRun() {
        return (boolean) ScenarioDataStore.get(DRY_RUN_MODE);
    }

    public static String getConfluenceUsernameFromScenarioDataStore() {
        return (String) ScenarioDataStore.get(CONFLUENCE_USERNAME);
    }

    public static String getConfluenceTokenFromScenarioDataStore() {
        return (String) ScenarioDataStore.get(CONFLUENCE_TOKEN);
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

    @AfterScenario
    public void setConfluenceUsernameAndTokenFromEnvVar() {
        ScenarioDataStore.remove(CONFLUENCE_USERNAME);
        ScenarioDataStore.remove(CONFLUENCE_TOKEN);
    }

    @AfterScenario(tags = {"create-space-manually"})
    public void afterScenario() {
        ConfluenceClient.deleteSpace(getScenarioSpaceKey());
    }

    @Step("Activate dry run mode")
    public void activateDryRunMode() {
        ScenarioDataStore.put(DRY_RUN_MODE, true);
    }

    @Step("Use Confluence user who does not have permission to create space")
    public void useConfluenceUserWhoDoesNotHavePermissionToCreateSpace() {
        ScenarioDataStore.put(CONFLUENCE_USERNAME, System.getenv("CONFLUENCE_USERNAME_WITHOUT_CREATE_SPACE"));
        ScenarioDataStore.put(CONFLUENCE_TOKEN, System.getenv("CONFLUENCE_TOKEN_WITHOUT_CREATE_SPACE"));
    }

    @Step("Space does not exist")
    public void assertSpaceDoesNotExist() {
        assertThat(ConfluenceClient.doesSpaceExist(getScenarioSpaceKey())).isFalse();
    }

    @Step("Space has name <name>")
    public void assertSpaceHasName(String name) {
        Space space = new Space(getScenarioSpaceKey());
        assertThat(space.getName()).isEqualTo(name);
    }

    @Step("Space has description <description>")
    public void assertSpaceHasDescription(String description) {
        Space space = new Space(getScenarioSpaceKey());
        assertThat(space.getDescription()).isEqualTo(description);
    }

    @Step("Published pages are: <table>")
    public void assertPublishedPages(Table expectedPages) throws Exception {
        int expectedTotal = expectedPages.getTableRows().size();
        assertConsoleSuccessOutput(expectedTotal - 1); // don't count the homepage as we don't publish it
        SpacePages spacePages = new SpacePages(getScenarioSpaceKey());
        assertThat(spacePages.total()).isEqualTo(expectedTotal);
        for (TableRow row : expectedPages.getTableRows()) {
            String actualParentPageTitle = spacePages.getParentPageTitle(row.getCell("title"));
            assertThat(actualParentPageTitle).isEqualTo(row.getCell("parent"));
        }
    }

    @Step("Specs <did|did not> get published")
    public void didPublishingOccur(String didPublishingOccur) throws IOException {
        boolean publishingOccurred = (didPublishingOccur.equalsIgnoreCase("did"));
        SpacePages spacePages = new SpacePages(getScenarioSpaceKey());
        if (publishingOccurred) {
            new Console().outputContains("Success: published");
            assertThat(spacePages.total()).isGreaterThan(1);
        } else {
            new Console().outputDoesNotContain("Success: published");
            assertThat(spacePages.total()).isEqualTo(1);
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
