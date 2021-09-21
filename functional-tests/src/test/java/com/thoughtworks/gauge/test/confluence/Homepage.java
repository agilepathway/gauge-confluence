package com.thoughtworks.gauge.test.confluence;

import org.json.JSONObject;

public class Homepage {

    private JSONObject jsonHomepage;

    public Homepage(JSONObject jsonHomepage) {
        this.jsonHomepage = jsonHomepage;
    }

    public String getTitle() {
        return jsonHomepage.getString("title");
    }

    public String getBody() {
        return jsonHomepage.getJSONObject("body").getJSONObject("view").getString("value");
    }

    @Override
    public String toString() {
        return jsonHomepage.toString(4);
    }
    
}
