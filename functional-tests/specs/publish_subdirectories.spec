# End to end example

## All specs are published, including those in subdirectories
The specs are published to Confluence as a page tree, mirroring the directory structure of the specs.  This means that
as well as a Confluence page being created for each spec, a page is also created for each directory.

* Publish specs to Confluence:

   |heading                              |path                        |
   |-------------------------------------|----------------------------|
   |specs dir spec                       |specs                       |
   |specs subfolder dir spec             |specs/subfolder             |
   |specs subfolder2 dir spec            |specs/subfolder2            |
   |specs subfolder subsubfolder dir spec|specs/subfolder/subsubfolder|

* Published pages are:

   |title                                |parent      |
   |-------------------------------------|------------|
   |Space Home                           |            |
   |specs                                |Space Home  |
   |specs dir spec                       |specs       |
   |subfolder                            |specs       |
   |specs subfolder dir spec             |subfolder   |
   |subfolder2                           |specs       |
   |specs subfolder2 dir spec            |subfolder2  |
   |subsubfolder                         |subfolder   |
   |specs subfolder subsubfolder dir spec|subsubfolder|
