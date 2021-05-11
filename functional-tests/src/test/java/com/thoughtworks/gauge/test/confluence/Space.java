package com.thoughtworks.gauge.test.confluence;

import org.json.JSONArray;
import org.json.JSONObject;

public class Space {

    private JSONArray pages;

    public Space(String key) {
        this.pages = ConfluenceClient.getAllPages(key);
    }

    public String getParentPageTitle(String title) {
        String parentPageID = getParentPageID(title);
        if (parentPageID.isEmpty())
            return "";
        JSONObject parentPage = getPageByID(parentPageID);
        return parentPage.getString("title");
    }

    private String getParentPageID(String title) {
        JSONObject page = getPageByTitle(title);
        JSONArray ancestors = page.getJSONArray("ancestors");
        if (ancestors.isEmpty())
            return "";
        JSONObject parent = ancestors.getJSONObject(ancestors.length() - 1);

        return parent.getString("id");
    }

    private JSONObject getPageByID(String pageID) {
        for (int i = 0; i < pages.length(); i++) {
            JSONObject page = pages.getJSONObject(i);
            if (page.getString("id").equals(pageID)) {
                return page;
            }
        }
        throw new IllegalStateException("page not found");
    }

    private JSONObject getPageByTitle(String title) {
        for (int i = 0; i < pages.length(); i++) {
            JSONObject page = pages.getJSONObject(i);
            if (page.getString("title").equals(title)) {
                return page;
            }
        }
        throw new IllegalStateException("page not found");
    }

    public int totalPages() {
        return pages.length();
    }

}
