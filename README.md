Gauge-Confluence
==========

[![Gauge Badge](https://gauge.org/Gauge_Badge.svg)](https://gauge.org)

[![build](https://github.com/agilepathway/gauge-confluence/workflows/build/badge.svg)](https://github.com/agilepathway/gauge-confluence/actions?query=workflow%3Abuild+event%3Apush+branch%3Amaster)
[![tests](https://github.com/agilepathway/gauge-confluence/workflows/FTs/badge.svg)](https://github.com/agilepathway/gauge-confluence/actions?query=workflow%3AFTs+event%3Apush+branch%3Amaster)
[![reviewdog](https://github.com/agilepathway/gauge-confluence/workflows/reviewdog/badge.svg)](https://github.com/agilepathway/gauge-confluence/actions?query=workflow%3Areviewdog+event%3Apush+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/agilepathway/gauge-confluence)](https://goreportcard.com/report/github.com/agilepathway/gauge-confluence)

[![releases](https://img.shields.io/github/v/release/agilepathway/gauge-confluence?color=blue&sort=semver)](https://github.com/agilepathway/gauge-confluence/releases)
[![License](https://img.shields.io/github/license/agilepathway/gauge-confluence?color=blue)](LICENSE)


Publishes Gauge specifications to Confluence. This is a plugin for [gauge](https://gauge.org/).
___
* [Why Publish Gauge Specs to Confluence](#why-publish-gauge-specs-to-confluence)
* [How to Use](#how-to-use)
  * [Typical Workflow](#typical-workflow)
  * [Supported Confluence versions](#supported-confluence-versions)
  * [Plugin setup](#plugin-setup)
  * [Running the plugin](#running-the-plugin-ie-publishing-specs-to-confluence)
  * [FAQs](#faqs)
* [Installation](#installation)
  * [Normal Installation](#normal-installation)
  * [Offline Installation](#offline-installation)
  * [Build from Source](#build-from-source)
* [Contributing](#contributing)
* [License](#license)

___


Why Publish Gauge Specs to Confluence
-------------------------------------

This plugin is aimed at teams who use Confluence for documentation. 

It enables [living documentation](https://www.infoq.com/articles/book-review-living-documentation/) by publishing 
your Gauge specs to Confluence and therefore allowing everyone to see them, seamlessly.

One of many use cases is if you are using [Specification by Example](http://specificationbyexample.com), for instance.

Writing your specifications in Gauge allows the documentation to be close to the code and also to be executable.
Being able to publish the specifications to the other tools you use makes this even more powerful.

See also the [Gauge Jira plugin](https://github.com/agilepathway/gauge-jira) which enables Gauge specs to be
published to Jira issues.  Using both of these plugins to publish your specs to Confluence and Jira is a powerful combo!

As [Gojko Adzic, the father of Specification by Example, says](https://gojko.net/2020/03/17/sbe-10-years.html#looking-forward-to-the-next-ten-years):

> *The big challenge related to tooling over the next 10 years will be in integrating better with Jira and its*
> *siblings. Somehow closing the loop so that teams that prefer to see information in task tracking tools get* 
> *the benefits of living documentation will be critical.*


How to Use
----------

### Typical Workflow

A typical workflow could be something like this:

1. collaborative story refinement sessions to come up with specification examples, using 
   [example mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) for instance
2. [write up the specification examples in Gauge](https://docs.gauge.org/writing-specifications.html)
3. use this plugin in a [Continuous Integration (CI) pipeline](https://www.thoughtworks.com/continuous-integration)
   to publish (or republish) the specifications to Confluence (and to Jira, if you use the Jira plugin too)
4. [automate the specifications using Gauge](https://docs.gauge.org/writing-specifications.html#step-implementations) 
   whenever possible (not essential, there's still value even when not automated)
5. continue the cycle throughout the lifespan of the story: more conversations, more spec updates, 
   more automated publishing to Confluence


### Supported Confluence versions

The plugin supports [Confluence Server](https://www.atlassian.com/software/confluence/download),
[Confluence Data Center](https://www.atlassian.com/software/confluence/download/data-center)
and [Confluence Cloud](https://www.atlassian.com/software/confluence).

If you find a problem with a particular version of Confluence, please
[raise an issue](../../issues)


### Plugin setup

There are three variables to configure, as either:

1. environment variables

2. properties in a 
   [properties file](https://docs.gauge.org/configuration.html#local-configuration-of-gauge-default-properties),
   e.g. `<project_root>/env/default/anythingyoulike.properties`

The three variables to configure are:

`CONFLUENCE_BASE_URL` e.g. `https://example.com` for Confluence Server, or `https://example.atlassian.net` for Confluence Cloud

`CONFLUENCE_USERNAME`

`CONFLUENCE_TOKEN`


### Running the plugin (i.e. publishing specs to Confluence)

`gauge docs confluence`

or, if you want to specify a different directory to the default `specs` directory

`gauge docs confluence <path to specs dir>`


### FAQs

1. Can the specifications be edited in Confluence and synced back into the Gauge specs?

   No.  We include a message in Confluence warning not to make edits to the specifications in Confluence.

2. Is it safe to publish the specs to Confluence multiple times?

   Yes.  The plugin replaces any previously published specs with the latest version.


Installation
------------


### Normal Installation

```
gauge install confluence
```
To install a specific version of the plugin use the ``--version`` flag.

```
gauge install confluence --version $VERSION
```


### Offline Installation

Download the plugin zip from the [Github Releases](https://github.com/agilepathway/gauge-confluence/releases),
or alternatively (if you want to experiment with an unreleased version, which is not recommended) from the
[artifacts](https://docs.github.com/actions/managing-workflow-runs/downloading-workflow-artifacts) in the
[`Store distros`](../../actions?query=workflow%3A%22Store+distros%22) GitHub Action (NB you must be logged
in to GitHub to be able to retrive the artifacts from there).

use the ``--file`` or ``-f`` flag to install the plugin from  zip file.

```
gauge install confluence --file ZIP_FILE_PATH
```

### Build from Source


#### Requirements
* [Golang](http://golang.org/)


#### Compiling

```
go run build/make.go
```

For cross-platform compilation

```
go run build/make.go --all-platforms
```


#### Installing
After compilation

```
go run build/make.go --install
```


#### Creating distributable

Note: Run after compiling

```
go run build/make.go --distro
```

For distributable across platforms: Windows and Linux for both x86 and x86_64

```
go run build/make.go --distro --all-platforms
```


Contributing
------------

See the [CONTRIBUTING.md](./CONTRIBUTING.md)


License
-------

`Gauge-Confluence` is released under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.