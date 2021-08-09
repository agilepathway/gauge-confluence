# Duplicate spec headings and directory names
The page title for every page in a Confluence space must be unique.
We use the specification heading as the Confluence page title.
And for directories we use the directory name as the Confluence page title.
So this means that there can be no duplicates among the spec headings and directory names.
NB the spec filenames are not relevant, it's just the spec headings and the directory names which must be unique.

This table shows examples with two specifications, as there needs to be at least two specs for there to be a
chance of duplicate spec headings.

   |spec 1 heading|spec 1 path  |spec 2 heading|spec 2 path      |result |message                                       |
   |--------------|-------------|--------------|-----------------|-------|----------------------------------------------|
   |one           |specs        |two           |specs            |Success|                                              |
   |one           |specs/folder |two           |specs/folder     |Success|                                              |
   |one           |specs/folder1|two           |specs/folder2    |Success|                                              |
   |same          |specs/folder1|same          |specs/folder2    |Failed |2 specs have the same heading                 |
   |one           |specs/same   |two           |specs/folder/same|Failed |2 directories have the same name              |
   |same          |specs/same   |two           |specs            |Failed |A spec heading and directory name are the same|


## Publishing to Confluence fails if there are any duplicate spec headings or directory names

* Publish specs to Confluence:

   |heading         |path         |
   |----------------|-------------|
   |<spec 1 heading>|<spec 1 path>|
   |<spec 2 heading>|<spec 2 path>|

* Output contains <result>
* Output contains <message>

## Republishing after fixing the duplicate spec headings or directory names works fine

* Publish specs to Confluence:

   |heading|
   |-------|
   |same   |
   |same   |

* Output contains "Failed: 2 specs have the same heading"

* Publish specs to Confluence:

   |heading      |
   |-------------|
   |same same    |
   |but different|

* Output contains "Success"
