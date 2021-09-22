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

    public int getVersion() {
        return jsonHomepage.getJSONObject("version").getInt("number");
    }

    @Override
    public String toString() {
        return jsonHomepage.toString(4);
    }
    
}
