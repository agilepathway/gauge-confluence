# Confluence Space creation

## If the Space does not already exist, the plugin will create it

* Publish "1" specs to Confluence

* Specs "did" get published

* Space has key "GITHUBCOMEXAMPLEUSEREXAMPLEREPO"

* Space has name "Gauge specs for example-user/example-repo"

The `example-user/example-repo` comes from the [dummy Git remote URL config in the test framework code][1].  When users run the plugin the Space name will be taken from the Git remote URL of the Git repository that the plugin is executed on.

* Output contains "Success: published 2 specs and directory pages to Confluence Space named: Gauge specs for example-user/example-repo"

* Space has description "Gauge (https://gauge.org) specifications from https://github.com/example-user/example-repo, published automatically by the Gauge Confluence plugin tool (https://github.com/agilepathway/gauge-confluence) as living documentation.  Do not edit this Space manually.  You can use Confluence's Include Macro (https://confluence.atlassian.com/doc/include-page-macro-139514.html) to include these specifications in as many of your existing Confluence Spaces as you wish."

## If the Space does not already exist, the plugin needs to be run by a user with permissions to create it
Tags: not-cloud
(Cannot run this test on Confluence Cloud as [Cloud's free tier does not provide configurable permissions][2])

* Use Confluence user who does not have permission to create space

* Space does not exist

* Publish "1" specs to Confluence and assert error: "rerun the plugin with a user who does have permissions to create the Space"

__________________________________________________________________________________________

[1]: https://github.com/agilepathway/gauge-confluence/blob/ae3129705ef85eaf56846d26b49e968de8b70e8b/functional-tests/src/test/java/com/thoughtworks/gauge/test/git/Config.java
[2]: https://support.atlassian.com/confluence-cloud/docs/manage-permissions-in-the-free-edition-of-confluence-cloud/