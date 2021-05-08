package com.thoughtworks.gauge.test.confluence;

import com.thoughtworks.gauge.Step;

public class SpaceManager {

    @Step("Delete space")
    public void deleteSpace() {
        ConfluenceClient.deleteSpace(spaceKey());
        System.out.println("Deleted space: " + spaceKey());
    }

    @Step("Create space")
    public void createSpace() {
        ConfluenceClient.createSpace(spaceKey(), spaceKey());
        System.out.println("Deleted space: " + spaceKey());
    }

    private String spaceKey() {
        return System.getenv("CONFLUENCE_SPACE_KEY");
    }

}
