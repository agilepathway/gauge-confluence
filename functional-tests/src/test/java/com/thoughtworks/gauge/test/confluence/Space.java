package com.thoughtworks.gauge.test.confluence;

import org.json.JSONObject;

public class Space {

    private JSONObject jsonSpace;

    public Space(String key) {
        this.jsonSpace = ConfluenceClient.getSpace(key);
    }

    public String getName() {
        return jsonSpace.getString("name");
    }

    @Override
    public String toString() {
        return jsonSpace.toString(4);
    }
    
}
