# End to end example
Tags: create-space-manually

## All specs are published, including those in subdirectories
The specs are published to Confluence as a page tree, mirroring the directory structure of the specs.  This means that
as well as a Confluence page being created for each spec, a page is also created for each directory.

* Publish specs to Confluence:

   |heading                       |path                        |
   |------------------------------|----------------------------|
   |A spec in the specs dir       |specs                       |
   |A spec in the subfolder dir   |specs/subfolder             |
   |A spec in the subfolder2 dir  |specs/subfolder2            |
   |A spec in the subsubfolder dir|specs/subfolder/subsubfolder|

* Published pages are:

   |title                         |parent      |
   |------------------------------|------------|
   |Space Home                    |            |
   |specs                         |Space Home  |
   |subfolder                     |specs       |
   |subfolder2                    |specs       |
   |A spec in the specs dir       |specs       |
   |subsubfolder                  |subfolder   |
   |A spec in the subfolder dir   |subfolder   |
   |A spec in the subfolder2 dir  |subfolder2  |
   |A spec in the subsubfolder dir|subsubfolder|
