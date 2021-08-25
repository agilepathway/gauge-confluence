# Confluence Space validity
This Gauge Confluence plugin requires that the specs for a given Gauge project are published to
their own dedicated [Confluence Space][1], with no manual edits or additions in Confluence.
Having the Space only contain published Gauge specs ensures that there is no danger of the plugin
inadvertently overwriting or deleting manually created Confluence pages.
NB The specs can still appear alongside your existing manually created Confluence documentation
in other spaces, by using Confluence's [Include Page macro][2].  This macro
[allows you to include the full page tree of specs in as many other spaces as you like][3].


## Publishing is aborted if the Space has been manually edited since the last publish

* Publish "1" specs to Confluence

* Manually add a page to the Confluence space

* Publish "1" specs to Confluence and assert error: "the space has been modified since the last publish"


## Publishing is aborted if the Space is not empty and has never been published to
The Space can have a homepage but no other pages before the first ever publish to the Space.

* Manually add a page to the Confluence space

* Publish "1" specs to Confluence and assert error: "the space must be empty when you publish for the first time"


## Publishing is aborted if the Space does not have a homepage

* Manually delete the Confluence space homepage

* Publish "1" specs to Confluence and assert error: "add a homepage manually in Confluence to the space, then try again"

__________________________________________________________________________________________

[1]: https://support.atlassian.com/confluence-cloud/docs/use-spaces-to-organize-your-work/
[2]: https://support.atlassian.com/confluence-cloud/docs/insert-the-include-page-macro/
[3]: https://community.atlassian.com/t5/Confluence-questions/How-do-create-cross-space-navigation-sidebar/qaq-p/441031
