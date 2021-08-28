# Duplicate spec headings and directory names
Tags: create-space-manually

The page title for every page in a Confluence space must be unique.
We use the specification heading as the Confluence page title.
And for directories we use the directory name as the Confluence page title.
So this means that there can be no duplicates among the spec headings and directory names.
NB the spec filenames are not relevant, it's just the spec headings and the directory names which must be unique.

This table shows examples with two specifications, as there needs to be at least two specs for there to be a
chance of duplicate spec headings.

   |spec 1 heading|spec 1 path  |spec 2 heading|spec 2 path      |did publishing occur?|message                                       |
   |--------------|-------------|--------------|-----------------|---------------------|----------------------------------------------|
   |one           |specs        |two           |specs            |did                  |                                              |
   |one           |specs/folder |two           |specs/folder     |did                  |                                              |
   |one           |specs/folder1|two           |specs/folder2    |did                  |                                              |
   |same          |specs/folder1|same          |specs/folder2    |did not              |2 specs have the same heading                 |
   |one           |specs/same   |two           |specs/folder/same|did not              |2 directories have the same name              |
   |same          |specs/same   |two           |specs            |did not              |A spec heading and directory name are the same|


## Publishing to Confluence fails if there are any duplicate spec headings or directory names

* Publish specs to Confluence and assert <did publishing occur?> succeed:

   |heading         |path         |
   |----------------|-------------|
   |<spec 1 heading>|<spec 1 path>|
   |<spec 2 heading>|<spec 2 path>|

* Specs <did publishing occur?> get published
* Output contains <message>

## Republishing after fixing the duplicate spec headings or directory names works fine

* Publish specs to Confluence and assert "did not" succeed:

   |heading|
   |-------|
   |same   |
   |same   |

* Specs "did not" get published

* Publish specs to Confluence and assert "did" succeed:

   |heading      |
   |-------------|
   |same same    |
   |but different|

* Specs "did" get published
