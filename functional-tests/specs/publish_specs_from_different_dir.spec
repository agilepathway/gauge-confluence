# Any directory containing specifications can be published

tags: create-space-manually

Gauge documentation plugins allow for any directory containing specs to be specified as a command-line argument to the plugin.

## Publish from a different directory than the default specs directory

* Publish "custom" directory to Confluence:

   |heading               |path  |
   |----------------------|------|
   |A spec in a custom dir|custom|

* Published pages are:

   |title                                         |parent                                        |
   |----------------------------------------------|----------------------------------------------|
   |Gauge specs for example-user/example-repo Home|                                              |
   |custom                                        |Gauge specs for example-user/example-repo Home|
   |A spec in a custom dir                        |custom                                        |
