package com.thoughtworks.gauge.test.confluence;

import com.thoughtworks.gauge.Step;

import org.json.JSONArray;
import org.json.JSONObject;

public class SpaceManager {

    @Step("Delete space")
    public void deleteSpace() {
        ConfluenceClient.deleteSpace(spaceKey());
        System.out.println("Deleted space: " + spaceKey());
    }

    @Step("Create space")
    public void createSpace() {
        ConfluenceClient.createSpace(spaceKey(), spaceKey());
        System.out.println("Created space: " + spaceKey());
    }

    @Step("Delete all spaces named <space name>")
    public void deleteAllSpacesNamed(String spaceName) {
        JSONArray spaces = ConfluenceClient.getAllSpaces();
        for (int i = 0; i < spaces.length(); i++) {
            deleteSpaceIfNamed(spaces.getJSONObject(i), spaceName);
        }
    }

    @Step("Print content for space with key <space key>")
    public void getContentForSpace(String spaceKey) {
        JSONArray content = ConfluenceClient.getAllPages(spaceKey);
        System.out.println(content);
    }

    private void deleteSpaceIfNamed(JSONObject sp, String spaceName) {
        if (spaceName.equals(sp.get("name"))) {
            ConfluenceClient.deleteSpace((String) sp.get("key"));
            System.out.println(String.format("Deleted space with key %s and name %s", sp.get("key"), sp.get("name")));
        }
    }

    private String spaceKey() {
        return System.getenv("CONFLUENCE_SPACE_KEY");
    }

}
