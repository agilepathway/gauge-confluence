# Dry run mode

Having a dry run mode is very useful, e.g. in a CI/CD pipeline the dry run mode can be run on feature branches and pull
requests to verify that the Gauge specs are in a valid state for publishing.  If they are not valid then the CI/CD 
pipeline build can fail, alerting the submitter of the pull request to amend them on the feature branch.  This ensures
that the Gauge specs are always in good shape to be automatically published by the CI/CD pipeline upon any push to the 
trunk branch (e.g. upon a successful pull request merge).

The dry run mode is set by setting a `DRY_RUN` [environment variable or property][2] with the value `true` (we can't
use a command-line flag for this as [Gauge does not propagate command line flags to documentation plugins][1]).


   |spec 1 heading|spec 2 heading|did publishing occur?|message                            |
   |--------------|--------------|---------------------|-----------------------------------|
   |same          |same          |did not              |Please fix the error then try again|
   |same          |different     |did not              |Dry run finished successfully      |


## Dry run mode indicates if specs are in a publishable state or not

* Activate dry run mode

* Publish specs to Confluence:

   |heading         |
   |----------------|
   |<spec 1 heading>|
   |<spec 2 heading>|

* Output contains <message>

* publishing <did publishing occur?> occur

__________________________________________________________________________________________

[1]: https://docs.gauge.org/configuration.html
[2]: https://github.com/getgauge/spectacle/issues/42#issuecomment-813483933
