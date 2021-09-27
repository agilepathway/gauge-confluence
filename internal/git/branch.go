package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const ciCommitBranchEnvVarName = "CI_COMMIT_BRANCH"

// discoverCurrentBranch returns the current Git branch name.
//
// The Git branch name is derived from the CI_COMMIT_BRANCH
// environment variable if the environment variable is present.
//
// This ensures that the Confluence plugin can discover the
// current branch properly when running on GitLab CI (as GitLab
// CI runs in detached HEAD mode when cloning or fetching from
// Git, meaning that the branch cannot be dynamically derived
// by the Confluence plugin).
//
// If the CI_COMMIT_BRANCH is not present then the branch is
// derived dynamically from the local git config, returning an
// error if it is not able to be derived.
func discoverCurrentBranch() (string, error) {
	ciCommitBranch, exists := os.LookupEnv(ciCommitBranchEnvVarName)
	if exists {
		return ciCommitBranch, nil
	}

	_, head, err := findHeadFile()
	if err != nil {
		return "", fmt.Errorf("there was a problem obtaining the HEAD file: %w", err)
	}

	currentBranch, err := discoverCurrentBranchFromHeadFile(head)
	if err != nil {
		return "", fmt.Errorf("there was a problem obtaining the current branch from the HEAD file: %w", err)
	}

	return currentBranch, nil
}

func findHeadFile() (string, string, error) {
	return findGitFile("HEAD")
}

func discoverCurrentBranchFromHeadFile(headPath string) (string, error) {
	if headPath == "" {
		return "", fmt.Errorf("no HEAD file defined")
	}

	headFile, err := ioutil.ReadFile(headPath) //nolint:gosec

	if err != nil {
		return "", fmt.Errorf("failed to load %s due to %s", headPath, err)
	}

	return abbreviatedHead(string(headFile))
}

func abbreviatedHead(fullHead string) (string, error) {
	if !strings.Contains(fullHead, "refs/heads/") {
		return "", fmt.Errorf("git is in detached HEAD state, HEAD is: %s", fullHead)
	}

	s := strings.TrimSpace(fullHead)
	s = strings.TrimPrefix(s, "ref:")
	s = strings.TrimSpace(s)

	return strings.TrimPrefix(s, "refs/heads/"), nil
}
