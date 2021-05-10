# Confluence Space validity

## Publishing is aborted if the Space has been manually edited since the last publish
This Gauge Confluence plugin requires that the specs for a given Gauge project are published to
their own dedicated [Confluence Space][1], with no manual edits or additions in Confluence.
Having the Space only contain published Gauge specs ensures that there is no danger of the plugin
inadvertently overwriting or deleting manually created Confluence pages.
NB The specs can still be appear alongside your existing manually created Confluence documentation
in other spaces, by using Confluence's [Include Page macro][2].  This macro
[allows you to include the full page tree of specs in as many other spaces as you like][3].

* Publish specs to Confluence:

   |heading|
   |-------|
   |A spec |

* Manually add a page to the Confluence space

* Publish specs to Confluence:

   |heading     |
   |------------|
   |Another spec|

* The error message "Failed: the space has been modified since the last publish" should be output


__________________________________________________________________________________________

[1]: https://support.atlassian.com/confluence-cloud/docs/use-spaces-to-organize-your-work/
[2]: https://support.atlassian.com/confluence-cloud/docs/insert-the-include-page-macro/
[3]: https://community.atlassian.com/t5/Confluence-questions/How-do-create-cross-space-navigation-sidebar/qaq-p/441031
