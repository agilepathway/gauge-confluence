# Concepts are not published

tags: create-space-manually

## Concepts are not published

* Publish specs to Confluence:

   |heading  |path |concept|
   |---------|-----|-------|
   |A spec   |specs|       |
   |A concept|specs|yes    |

* Published pages are:

   |title                                         |parent                                        |
   |----------------------------------------------|----------------------------------------------|
   |Gauge specs for example-user/example-repo Home|                                              |
   |specs                                         |Gauge specs for example-user/example-repo Home|
   |A spec                                        |specs                                         |


## A directory that just contains concepts is not published

* Publish specs to Confluence:

   |heading                    |path          |concept|
   |---------------------------|--------------|-------|
   |A spec in a specs dir      |specs         |       |
   |A concept in a concepts dir|specs/concepts|yes    |

* Published pages are:

   |title                                         |parent                                        |
   |----------------------------------------------|----------------------------------------------|
   |Gauge specs for example-user/example-repo Home|                                              |
   |specs                                         |Gauge specs for example-user/example-repo Home|
   |A spec in a specs dir                         |specs                                         |


## Nested directories that just contain concepts are not published

* Publish specs to Confluence:

   |heading                                  |path                   |concept|
   |-----------------------------------------|-----------------------|-------|
   |A spec in a specs dir                    |specs                  |       |
   |A concept in a concepts dir              |specs/concepts         |yes    |
   |A concept in a sub concepts dir          |specs/concepts/sub     |yes    |
   |Another concept in a sub concepts dir    |specs/concepts/sub     |yes    |
   |Another concept in a sub sub concepts dir|specs/concepts/sub/sub2|yes    |

* Published pages are:

   |title                                         |parent                                        |
   |----------------------------------------------|----------------------------------------------|
   |Gauge specs for example-user/example-repo Home|                                              |
   |specs                                         |Gauge specs for example-user/example-repo Home|
   |A spec in a specs dir                         |specs                                         |
