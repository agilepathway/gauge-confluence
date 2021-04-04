package com.thoughtworks.gauge.test.confluence;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.HttpRequest.BodyPublishers;
import java.util.Base64;

import org.apache.commons.lang3.StringUtils;
import org.json.JSONArray;
import org.json.JSONObject;

public class ConfluenceClient {

    public static void createSpace(String spaceKey, String spaceName) {
        sendConfluenceRequest(createSpaceRequest(spaceKey, spaceName));
    }

    public static void deleteSpace(String spaceKey) {
        sendConfluenceRequest(deleteSpaceRequest(spaceKey));
    }

    public static JSONArray getAllPages(String spaceKey) {
        HttpResponse<String> rawResponse = sendConfluenceRequest(getAllPagesRequest(spaceKey));
        JSONObject jsonResponse = new JSONObject(rawResponse.body());
        return (JSONArray) jsonResponse.get("results");
    }

    private static HttpRequest createSpaceRequest(String spaceKey, String spaceName) {
        JSONObject description = new JSONObject().put("plain",
                new JSONObject().put("value", spaceName).put("representation", "plain"));
        JSONObject body = new JSONObject().put("key", spaceKey).put("name", spaceName).put("description", description);
        HttpRequest.Builder builder = baseConfluenceRequest();
        builder.uri(URI.create(baseSpaceAPIURL()));
        builder.POST(BodyPublishers.ofString(body.toString()));
        return builder.build();
    }

    private static HttpRequest deleteSpaceRequest(String spaceKey) {
        HttpRequest.Builder builder = baseConfluenceRequest();
        String deleteSpaceURL = String.format("%1$s/%2$s", baseSpaceAPIURL(), spaceKey);
        builder.uri(URI.create(deleteSpaceURL));
        builder.DELETE();
        return builder.build();
    }

    private static HttpRequest getAllPagesRequest(String spaceKey) {
        HttpRequest.Builder builder = baseConfluenceRequest();
        String getAllPagesURL = String.format("%1$s?spaceKey=%2$s&expand=ancestors", baseContentAPIURL(), spaceKey);
        builder.uri(URI.create(getAllPagesURL));
        return builder.build();
    }

    private static HttpResponse<String> sendConfluenceRequest(HttpRequest request) {
        HttpClient client = HttpClient.newBuilder().version(HttpClient.Version.HTTP_1_1).build();
        try {
            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            verifySuccessfulResponse(response);
            return response;
        } catch (IOException | InterruptedException e) {
            throw new IllegalStateException("Exception when sending Confluence space request", e);
        }
    }

    private static void verifySuccessfulResponse(HttpResponse<String> response) {
        if (response.statusCode() > 299) {
            String message = String.format("Expected 2xx response but got %1$s. Response body: %2$s",
                    response.statusCode(), response.body());
            throw new IllegalStateException(message);
        }
    }

    private static HttpRequest.Builder baseConfluenceRequest() {
        String confluenceUsername = System.getenv("CONFLUENCE_USERNAME");
        String confluenceToken = System.getenv("CONFLUENCE_TOKEN");
        return HttpRequest.newBuilder().header("Content-Type", "application/json").header("Authorization",
                basicAuth(confluenceUsername, confluenceToken));
    }

    private static String baseSpaceAPIURL() {
        return String.format("%1$s/rest/api/space", confluenceBaseURL());
    }

    private static String baseContentAPIURL() {
        return String.format("%1$s/rest/api/content", confluenceBaseURL());
    }

    private static String basicAuth(String username, String password) {
        return "Basic " + Base64.getEncoder().encodeToString((username + ":" + password).getBytes());
    }

    private static String confluenceBaseURL() {
        String baseURL = System.getenv("CONFLUENCE_BASE_URL");
        if (baseURL.endsWith("/")) {
            baseURL = StringUtils.chop(baseURL);
        }
        if (baseURL.endsWith("atlassian.net")) {
            return String.format("%1$s/wiki", baseURL);
        }
        return baseURL;
    }

}
