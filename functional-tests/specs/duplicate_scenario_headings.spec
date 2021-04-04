# Duplicate spec headings and directory names
The page title for every page in a Confluence space must be unique.
We use the specification heading as the Confluence page title.
And for directories we use the directory name as the Confluence page title.
So this means that there can be no duplicates among the spec headings and directory names.
NB the spec filenames are not relevant, it's just the spec headings and the directory names which must be unique.

This table shows examples with two specifications, as there needs to be at least two specs for there to be a
chance of duplicate spec headings.

   |spec 1 heading|spec 1 path     |spec 2 heading|spec 2 path         |result |message             |
   |--------------|----------------|--------------|--------------------|-------|--------------------|
   |one           |specs           |two           |specs               |Success|                    |
   |one           |specs/folder    |two           |specs/folder        |Success|                    |
   |one           |specs/subfolder1|two           |specs/subfolder2    |Success|                    |
   |same          |specs/subfolder1|same          |specs/subfolder2    |Failed |duplicate page: same|
   |one           |specs/same      |two           |specs/subfolder/same|Failed |duplicate page: same|
   |same          |specs/same      |two           |specs               |Failed |duplicate page: same|


## Publishing to Confluence fails if there are any duplicate spec headings or directory names

* Publish specs to Confluence:

   |heading         |path         |
   |----------------|-------------|
   |<spec 1 heading>|<spec 1 path>|
   |<spec 2 heading>|<spec 2 path>|

* Output contains <result>
* Output contains <message>

