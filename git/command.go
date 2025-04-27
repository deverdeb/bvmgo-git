package git

import (
	"fmt"
	"github.com/deverdeb/bvmgo-git/utils/caches"
	"os/exec"
	"strings"
	"time"
)

var branchesCache = caches.BuildCache[string, *Branch]()
var tagsCache = caches.BuildCache[string, *Tag]()
var commitsCache = caches.BuildCache[string, *Commit]()

func GetHead() (Hash, error) {
	commandResult, err := executeCommand("git", "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("faild to read Git HEAD hash: %w", err)
	}
	if len(commandResult) == 0 {
		return "", fmt.Errorf("faild to read Git HEAD hash: HEAD hash not found")
	}
	return Hash(commandResult), nil
}

func ListBranches() ([]*Branch, error) {
	result := make([]*Branch, 0)
	commandResult, err := executeCommand("git", "for-each-ref", "--format=%(refname)", "refs/heads", "refs/remotes")
	if err != nil {
		return result, fmt.Errorf("failed to read Git branches: %w", err)
	}
	lines := strings.Split(commandResult, "\n")
	for _, line := range lines {
		branchName := strings.TrimSpace(line)
		if len(branchName) > 0 {
			branch, err := GetBranch(branchName)
			if err != nil {
				return result, err
			}
			if branch != nil {
				result = append(result, branch)
			}
		}
	}
	return result, nil
}

func GetBranch(branchName string) (*Branch, error) {
	branch := branchesCache.Get(branchName)
	if branch != nil {
		return branch, nil
	}
	commit, err := GetInfoCommitByName(branchName)
	if err != nil {
		return nil, fmt.Errorf("failed to read Git commit for branch %s: %w", branchName, err)
	}
	branch = &Branch{
		Name:   branchName,
		Commit: commit,
	}
	branchesCache.Put(branchName, branch)
	return branch, nil
}

func ListTags() ([]*Tag, error) {
	result := make([]*Tag, 0)
	commandResult, err := executeCommand("git", "for-each-ref", "--format=<TAGINFO>|%(refname)|%(taggerdate:format:%Y-%m-%dT%H:%M:%S)|%(taggername)|%(subject)", "refs/tags")
	if err != nil {
		return result, fmt.Errorf("faild to read Git tags: %w", err)
	}
	lines := strings.Split(commandResult, "<TAGINFO>|")
	for _, line := range lines {
		cleanLine := strings.TrimSpace(line)
		if len(cleanLine) == 0 {
			continue
		}
		tagInfo, err := parseCommandLineResultToInfoTag(cleanLine)
		if err != nil {
			return result, fmt.Errorf("faild to extract tag information: %w", err)
		}
		if tagInfo != nil {
			tagsCache.Put(tagInfo.Name, tagInfo)
			result = append(result, tagInfo)
		}
	}
	return result, nil
}

func GetTag(tagName string) (*Tag, error) {
	tag := tagsCache.Get(tagName)
	if tag != nil {
		return tag, nil
	}
	if tagsCache.Size() == 0 {
		_, err := ListTags()
		if err != nil {
			return nil, err
		}
		tag = tagsCache.Get(tagName)
	}
	return tag, nil
}

func parseCommandLineResultToInfoTag(lineCommandResult string) (*Tag, error) {
	fields := strings.SplitN(lineCommandResult, "|", 4)
	if fields == nil || len(fields) < 4 {
		return nil, fmt.Errorf("failed to read Git commit information: Missing information - Command result: %s", lineCommandResult)
	}
	var err error = nil
	var when = time.Time{}
	if len(fields[2]) > 0 {
		when, err = time.Parse("2006-01-02T15:04:05", fields[1])
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag date: %w - Command result: %s", err, lineCommandResult)
	}
	tagName := strings.TrimSpace(fields[0])
	commit, err := GetInfoCommitByName(tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to read Git commit for tag %s: %w", tagName, err)
	}
	infoTag := Tag{
		Name:    strings.TrimSpace(fields[0]),
		When:    when,
		Author:  strings.TrimSpace(fields[2]),
		Subject: strings.TrimSpace(fields[3]),
		Commit:  commit,
	}
	return &infoTag, nil
}

func GetInfoCommitByHash(hash Hash) (*Commit, error) {
	return GetInfoCommitByName(hash.String())
}

func GetInfoCommitByName(branchOrTagName string) (*Commit, error) {
	// https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
	if len(branchOrTagName) == 0 {
		return nil, fmt.Errorf("failed to read Git commit from empty hash")
	}
	var err error = nil
	commit := commitsCache.GetOrCompute(branchOrTagName, func(branchOrTagName string) *Commit {
		var found *Commit = nil
		found, err = searchCommitToGit(branchOrTagName)
		if err != nil {
			return nil
		} else {
			return found
		}
	})
	return commit, err
}

func searchCommitToGit(branchOrTagName string) (*Commit, error) {
	// https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
	commandResult, err := executeCommand("git", "log", "-1", "--date=format:%Y-%m-%dT%H:%M:%S", "--pretty=format:%H|%P|%ad|%an|%s", branchOrTagName)
	if err != nil {
		return nil, fmt.Errorf("failed to read Git commit information: %w", err)
	}
	return parseCommandResultToInfoCommit(commandResult)
}

func parseCommandResultToInfoCommit(commandResult string) (*Commit, error) {
	fields := strings.SplitN(commandResult, "|", 5)
	if fields == nil || len(fields) < 5 {
		return nil, fmt.Errorf("failed to read Git commit information: Missing information - Command result: %s", commandResult)
	}
	when, err := time.Parse("2006-01-02T15:04:05", fields[2])
	if err != nil {
		return nil, fmt.Errorf("failed to read commit date: %w - Command result: %s", err, commandResult)
	}
	parents := strings.Split(fields[1], " ")
	parentHashes := make([]Hash, len(parents))
	for idx, parent := range strings.Split(fields[1], " ") {
		parentHashes[idx] = Hash(strings.TrimSpace(parent))
	}
	infoCommit := Commit{
		Hash:    Hash(strings.TrimSpace(fields[0])),
		Parents: parentHashes,
		When:    when,
		Author:  strings.TrimSpace(fields[3]),
		Subject: strings.TrimSpace(fields[4]),
	}
	return &infoCommit, nil
}

func executeCommand(application string, arguments ...string) (string, error) {
	command := exec.Command(application, arguments...)
	stdout, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command %s %v: %w\n", application, arguments, err)
	}
	return string(stdout), nil
}
