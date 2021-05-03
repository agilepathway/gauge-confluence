package confluence

import "fmt"

const pluginAbortedMessage = "The plugin aborted as soon as it encountered the duplicate.\n"

type duplicatePageError struct {
	published page
	duplicate page
}

func (e *duplicatePageError) Error() string {
	var msg string

	switch {
	case e.published.isSpec() && e.duplicate.isSpec():
		msg = e.duplicateSpecHeadings()
	case e.published.isDir && e.duplicate.isDir:
		msg = e.duplicateDirNames()
	case e.published.isSpec() && e.duplicate.isDir:
		msg = e.duplicateSpecHeadingAndDirName(e.published, e.duplicate)
	case e.published.isDir && e.duplicate.isSpec():
		msg = e.duplicateSpecHeadingAndDirName(e.duplicate, e.published)
	}

	return msg + `(This is because the page title for every page in a Confluence space must be unique.
We use the specification heading as the Confluence page title.
And for directories we use the directory name as the Confluence page title.
So this means that there can be no duplicates among the spec headings and directory names.
NB the spec filenames are not relevant, it's just the spec headings and the directory names which must be unique.)`
}

func (e *duplicatePageError) duplicateSpecHeadings() string {
	return fmt.Sprintf("2 specs have the same heading: \"%s\"\n", e.duplicate.title) +
		fmt.Sprintf("The paths for the 2 specs containing the duplicate heading: \n%s\n%s\n",
			e.published.path, e.duplicate.path) +
		pluginAbortedMessage +
		"Change one of the spec headings and then run the Gauge Confluence plugin again.\n"
}

func (e *duplicatePageError) duplicateDirNames() string {
	return fmt.Sprintf("2 directories have the same name: \"%s\"\n", e.duplicate.title) +
		fmt.Sprintf("The 2 directory paths containing the duplicate name: %s %s\n",
			e.published.path, e.duplicate.path) +
		pluginAbortedMessage +
		"Change one of the directory names and then run the Gauge Confluence plugin again.\n"
}

func (e *duplicatePageError) duplicateSpecHeadingAndDirName(specPage, dirPage page) string {
	return fmt.Sprintf("A spec heading and directory name are the same: \"%s\"\n", specPage.title) +
		fmt.Sprintf("Spec path (for the spec with the heading which clashes with the directory name): %s\n", specPage.path) +
		fmt.Sprintf("Directory path (with the directory name which clashes with the spec heading): %s\n", dirPage.path) +
		pluginAbortedMessage +
		"Change the spec heading or the directory name and then run the Gauge Confluence plugin again.\n"
}
