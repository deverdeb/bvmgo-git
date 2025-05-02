package git

import (
	"fmt"
	"regexp"
	"time"
)

var ignoreBefore = time.Date(
	2020, time.January, 1, 00, 00, 00, 000000000, time.UTC)

var branchesRegex = regexp.MustCompile("^refs/(heads|remotes/origin)/(.*)")
var tagsRegex = regexp.MustCompile("^refs/tags/(.*)")

type extractionContext struct {
	hashList ConsolidatedCommitMap
}

func ExtractCommits() (ConsolidatedCommitMap, error) {
	var hashList = NewConsolidatedCommitMap()

	context := extractionContext{
		hashList: hashList,
	}

	err := extractBranches(context)
	if err != nil {
		return nil, fmt.Errorf("could not extract git branches information\n - cause by: %w", err)
	}

	err = extractTags(context)
	if err != nil {
		return nil, fmt.Errorf("could not extract git tags information\n - cause by: %w", err)
	}

	return hashList, nil
}

func extractBranches(context extractionContext) error {
	branches, err := ListBranches()
	if err != nil {
		return fmt.Errorf("could not read git branches\n - cause by: %w", err)
	}

	fmt.Printf("\nbranches:\n")
	for _, branch := range branches {
		if !branchesRegex.MatchString(branch.Name) {
			return nil
		}
		fmt.Printf(" - branch %s\n", branch)
		hashInfo := findConsolidatedCommitByName(context, branch.Name)
		if hashInfo != nil {
			hashInfo.Branches = append(hashInfo.Branches, branch.Name)
		}
	}
	return nil
}

func extractTags(context extractionContext) error {
	tags, err := ListTags()
	if err != nil {
		return fmt.Errorf("could not read git tags\n - cause by: %w", err)
	}

	fmt.Printf("\ntags:\n")
	for _, tag := range tags {
		if !tagsRegex.MatchString(tag.Name) {
			return nil
		}

		fmt.Printf(" - tag %s\n", tag)
		hashInfo := findConsolidatedCommitByName(context, tag.Name)
		if hashInfo != nil {
			hashInfo.Tags = append(hashInfo.Tags, tag)
		}
	}
	return nil
}

func findConsolidatedCommitByHash(context extractionContext, hash Hash) *ConsolidatedCommit {
	if len(hash) == 0 {
		return nil
	}
	commitInfo := context.hashList[hash]
	if commitInfo != nil {
		return commitInfo
	}
	commit, err := GetInfoCommitByHash(hash)
	if err != nil {
		fmt.Printf("   (!) > failed to read information for hash [%s]: %s\n", hash.String(), err.Error())
		return nil
	}
	commitInfo = buildConsolidatedCommitFromCommit(context, commit)
	if commitInfo != nil {
		context.hashList[commit.Hash] = commitInfo
	}
	return commitInfo

}

func findConsolidatedCommitByName(context extractionContext, branchOrTagName string) *ConsolidatedCommit {
	if len(branchOrTagName) == 0 {
		return nil
	}
	commit, err := GetInfoCommitByName(branchOrTagName)
	if err != nil {
		fmt.Printf("   (!) > failed to read information for [%s]: %s\n", branchOrTagName, err.Error())
		return nil
	}
	commitInfo := context.hashList[commit.Hash]
	if commitInfo != nil {
		return commitInfo
	}
	commitInfo = buildConsolidatedCommitFromCommit(context, commit)
	if commitInfo != nil {
		context.hashList[commit.Hash] = commitInfo
	}
	return commitInfo
}

func buildConsolidatedCommitFromCommit(context extractionContext, commit *Commit) *ConsolidatedCommit {
	when := commit.When
	if ignoreBefore.After(when) {
		// Ignore to old commit
		return nil
	}
	fmt.Printf("   > process commit %s - date %s\n", string(commit.Hash), when.Format("02/01/2006 15:04:05"))
	commitInfo := BuildConsolidatedCommit(commit)
	for _, parentHash := range commit.Parents {
		parentCommit := findConsolidatedCommitByHash(context, parentHash)
		if parentCommit != nil {
			commitInfo.HashParents.Add(parentCommit.Commit.Hash)
			parentCommit.HashChildren.Add(commitInfo.Commit.Hash)
		}
	}
	return commitInfo
}
