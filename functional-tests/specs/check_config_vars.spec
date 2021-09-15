# Configuration variables are required to be set, unless in dry run mode
Tags: create-space-manually

   |config variable     |
   |--------------------|
   |CONFLUENCE_BASE_URL |
   |CONFLUENCE_USERNAME |
   |CONFLUENCE_TOKEN    |

## The plugin fails if required configuration variables are not set

* Required configuration variable <config variable> must be set


## Configuration variables are not required to be set in dry run mode

* Activate dry run mode

* Initialize an empty Gauge project

* Publish Confluence Documentation for the current project with no <config variable> configured and assert "did" succeed

* Output contains "Dry run finished successfully"


________________________________________________________________________________________________

Read more about [Gauge's configuration documentation](https://docs.gauge.org/configuration.html)
