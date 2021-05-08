# Spec validity for publishing to Confluence

## Specs without a heading are skipped

* Publish specs to Confluence:

   |heading|
   |-------|
   |spec 1 |
   |       |
   |spec 3 |

* Published pages are:

   |title                              |parent                             |
   |-----------------------------------|-----------------------------------|
   |Temporary Gauge Scenario Space Home|                                   |
   |specs                              |Temporary Gauge Scenario Space Home|
   |spec 1                             |specs                              |
   |spec 3                             |specs                              |

* Output contains "Skipping file: could not find a spec heading"
