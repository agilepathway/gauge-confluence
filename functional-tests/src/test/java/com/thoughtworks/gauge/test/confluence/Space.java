package com.thoughtworks.gauge.test.confluence;

import org.json.JSONObject;

public class Space {

    private JSONObject jsonSpace;

    public Space(String key) {
        this.jsonSpace = ConfluenceClient.getSpace(key);
    }

    public String getKey() {
        return jsonSpace.getString("key");
    }

    public String getName() {
        return jsonSpace.getString("name");
    }

    public String getDescription() {
        return jsonSpace.getJSONObject("description").getJSONObject("plain").getString("value");
    }

    public Homepage getHomepage() {
        return new Homepage(jsonSpace.getJSONObject("homepage"));
    }

    @Override
    public String toString() {
        return jsonSpace.toString(4);
    }
    
}
