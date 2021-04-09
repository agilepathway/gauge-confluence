package com.thoughtworks.gauge.test.confluence;

import com.thoughtworks.gauge.BeforeScenario;
import com.thoughtworks.gauge.AfterScenario;
import com.thoughtworks.gauge.datastore.ScenarioDataStore;

import java.time.Instant;

public class Confluence {

    private static final String SCENARIO_SPACE_KEY_NAME = "confluence-space-key";
    private static final String SCENARIO_SPACE_NAME = "Temporary Gauge Scenario Space";

    public static String getScenarioSpaceKey() {
        return (String) ScenarioDataStore.get(SCENARIO_SPACE_KEY_NAME);
    }

    @BeforeScenario
    public void BeforeScenario() {
        ScenarioDataStore.put(SCENARIO_SPACE_KEY_NAME, currentTimeInMilliseconds());
        ConfluenceClient.createSpace(getScenarioSpaceKey(), SCENARIO_SPACE_NAME);
    }

    @AfterScenario
    public void AfterScenario() {
        ConfluenceClient.deleteSpace(getScenarioSpaceKey());
    }

    public String currentTimeInMilliseconds() {
        return String.valueOf(Instant.now().toEpochMilli());
    }

}
