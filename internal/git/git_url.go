package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// copied from https://github.com/go-git/go-git/blob/bf3471db54b0255ab5b159005069f37528a151b7/internal/url/url.go#L9
var scpLikeURIRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})(?:\/|:))?(?P<path>[^\\].*\/[^\\].*)$`) //nolint:lll

// WebURL provides the web URL of the remote Git repository.
func WebURL() (*url.URL, error) {
	remoteGitURL, err := discoverRemoteGitURL()

	if err != nil {
		return nil, err
	}

	return buildGitWebURL(remoteGitURL)
}

// RemoteURLPath provides the path part of the remote Git URL.
// It does not include the slash at the beginning of the path.
// e.g. user-or-org/project
func RemoteURLPath() (string, error) {
	gitWebURL, err := WebURL()
	if err != nil {
		return "", err
	}

	return parseURLPath(gitWebURL)
}

// SpecGitURL gives the remote Git URL (e.g. on GitHub, GitLab, Bitbucket etc) for a spec file
func SpecGitURL(absoluteSpecPath, projectRoot string) string {
	gitWebURL, err := WebURL()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	relativeSpecPath := strings.TrimPrefix(absoluteSpecPath, projectRoot)

	branch, err := discoverCurrentBranch()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return gitWebURL.String() + "/blob/" + branch + toURLFormat(relativeSpecPath)
}

func parseURLPath(u *url.URL) (string, error) {
	return strings.TrimPrefix(u.Path, "/"), nil
}

// buildGitWebURL constructs the publicly accessible Git web URL from a Git remote URL.
func buildGitWebURL(remoteGitURI string) (*url.URL, error) {
	url, err := url.Parse(remoteGitURI)

	isStandardURL := err == nil && url != nil
	if isStandardURL {
		webURL := gitWebURLScheme(url.Scheme) + "://" + url.Host + strings.TrimSuffix(url.Path, ".git")
		return url.Parse(webURL)
	}

	if isSCPStyleURI(remoteGitURI) {
		_, host, port, path := findScpLikeComponents(remoteGitURI)
		webURL := "https://" + hostAndPort(host, port) + "/" + strings.TrimSuffix(path, ".git")

		return url.Parse(webURL)
	}

	return nil, fmt.Errorf("could not parse Git URL %s", remoteGitURI)
}

func hostAndPort(host, port string) string {
	if port == "" {
		return host
	}

	return host + ":" + port
}

func isSCPStyleURI(input string) bool {
	return scpLikeURIRegExp.MatchString(input)
}

func findScpLikeComponents(uri string) (user, host, port, path string) {
	m := scpLikeURIRegExp.FindStringSubmatch(uri)
	return m[1], m[2], m[3], m[4]
}

func gitWebURLScheme(input string) string {
	if input == "http" {
		return input
	}

	return "https"
}

// toURLFormat converts any Windows path slashes to URL format (i.e. forward slashes).
func toURLFormat(input string) string {
	return strings.ReplaceAll(input, "\\", "/")
}
