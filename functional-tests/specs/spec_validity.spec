# Spec validity for publishing to Confluence

tags: create-space-manually

## Specs without a heading are not published to Confluence

* Publish specs to Confluence:

   |has heading|heading|
   |-----------|-------|
   |yes        |Spec 1 |
   |no         |       |
   |yes        |Spec 3 |

* Published pages are:

   |title                                         |parent                                        |
   |----------------------------------------------|----------------------------------------------|
   |Gauge specs for example-user/example-repo Home|                                              |
   |specs                                         |Gauge specs for example-user/example-repo Home|
   |Spec 1                                        |specs                                         |
   |Spec 3                                        |specs                                         |

* Output contains "Skipping file: could not find a spec heading"
