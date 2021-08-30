# Dry run mode

Having a dry run mode is very useful, e.g. in a CI/CD pipeline the dry run mode can be run on feature branches and pull
requests to verify that the Gauge specs are in a valid state for publishing.  If they are not valid then the CI/CD
pipeline build can fail, alerting the submitter of the pull request to amend them on the feature branch.  This ensures
that the Gauge specs are always in good shape to be automatically published by the CI/CD pipeline upon any push to the
trunk branch (e.g. upon a successful pull request merge).

The Gauge Confluence plugin will exit with a non-zero exit code if the dry run finds that any specs are not in a valid
state for publishing.  This means that the CI/CD pipeline can use the exit code (0 for success, not 0 for fail) to pass
or fail the dry run build (and indeed the actual build too, when not running in dry run mode).

The dry run mode is set by setting a `DRY_RUN` [environment variable or property][2] with the value `true` (we can't
use a command-line flag for this as [Gauge does not propagate command line flags to documentation plugins][1]).


   |spec 1 heading|spec 2 heading|did dry run succeed?|message                      |
   |--------------|--------------|--------------------|-----------------------------|
   |same          |same          |did not             |Failed to generate docs      |
   |same          |different     |did                 |Dry run finished successfully|


## Dry run mode indicates if specs are in a publishable state or not
Tags: create-space-manually

* Activate dry run mode

* Publish specs to Confluence and assert <did dry run succeed?> succeed:

   |heading         |
   |----------------|
   |<spec 1 heading>|
   |<spec 2 heading>|

* Output contains <message>

* Specs "did not" get published


## Dry run mode does not require the Confluence Space to exist yet
The absence of the "create-space-manually" tag means the Confluence Space does not
exist for this scenario

* Space does not exist

* Activate dry run mode

* Publish "1" specs to Confluence

* Output contains "Dry run finished successfully"

* Space does not exist


__________________________________________________________________________________________

[1]: https://docs.gauge.org/configuration.html
[2]: https://github.com/getgauge/spectacle/issues/42#issuecomment-813483933
