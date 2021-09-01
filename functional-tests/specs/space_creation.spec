# Confluence Space creation

## If the Space does not already exist, the plugin will create it

* Space does not exist

* Publish "1" specs to Confluence

* Specs "did" get published

* Space has name "Gauge specs for example-user/example-repo"

The `example-user/example-repo` comes from the [dummy Git remote URL config in the
test framework code][1].  When users run the plugin the Space name will be taken from 
the Git remote URL of the Git repository that the plugin is executed on.


__________________________________________________________________________________________

[1]: https://github.com/agilepathway/gauge-confluence/blob/ae3129705ef85eaf56846d26b49e968de8b70e8b/functional-tests/src/test/java/com/thoughtworks/gauge/test/git/Config.java
