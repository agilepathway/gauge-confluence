# Delete existing published specs before publishing

## The plugin deletes all existing published specs before publishing
It is safe to do this as before doing so we abort if the Space has been manually edited since the last publish.
NB The Space homepage is not deleted, it is just the published specs that are deleted (i.e. all the other pages
in the Space apart from the homepage).

* Publish "26" specs to Confluence
The default limit for pagination on the Confluence API is 25, so we publish 26 specs here to ensure that the
plugin handles pagination correctly when deleting the specs on the second publish below.

* Publish specs to Confluence:

   |heading     |
   |------------|
   |Another spec|

* Published pages are:

   |title                         |parent      |
   |------------------------------|------------|
   |Space Home                    |            |
   |specs                         |Space Home  |
   |Another spec                  |specs       |
