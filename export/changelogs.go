package export

import (
	"fmt"
	"github.com/deverdeb/bvmgo-git/git"
	"os"
	"strings"
)

func ExportChangelogs(commitsMap map[git.Hash]*git.ConsolidatedCommit) error {
	commitsWithTag := extractCommitsWithTags(commitsMap)
	changelogsByTag := extractChangelogsFromTags(commitsWithTag, commitsMap)
	return writeChangelogs(changelogsByTag)
}

func extractCommitsWithTags(commitsMap map[git.Hash]*git.ConsolidatedCommit) []*git.ConsolidatedCommit {
	commitsWithTag := make([]*git.ConsolidatedCommit, 0)
	for _, commitInfo := range commitsMap {
		if commitInfo.IsTag() {
			commitsWithTag = append(commitsWithTag, commitInfo)
		}
	}
	return commitsWithTag
}

func extractChangelogsFromTags(commitsWithTag []*git.ConsolidatedCommit, commitsMap map[git.Hash]*git.ConsolidatedCommit) map[*git.Tag][]*git.ConsolidatedCommit {
	changelogs := make(map[*git.Tag][]*git.ConsolidatedCommit)
	for _, commitWithTag := range commitsWithTag {
		for _, tag := range commitWithTag.Tags {
			changelogs[tag] = extractChangelogFromTag(commitWithTag.Commit.Hash, commitsMap)
		}
	}
	return changelogs
}

func extractChangelogFromTag(tagHash git.Hash, commitsMap map[git.Hash]*git.ConsolidatedCommit) []*git.ConsolidatedCommit {
	return completeTagCommitsForChangelog(make([]*git.ConsolidatedCommit, 0), tagHash, commitsMap)
}

func completeTagCommitsForChangelog(changelog []*git.ConsolidatedCommit, hash git.Hash, commitsMap map[git.Hash]*git.ConsolidatedCommit) []*git.ConsolidatedCommit {
	commitInfo := commitsMap[hash]
	if commitInfo == nil {
		return changelog
	}
	if len(changelog) != 0 && commitInfo.IsTag() {
		// Ce n'est pas le premier (le tag de départ) et c'est un tag
		// => C'est la fin du changelog
		return changelog
	}
	changelog = append(changelog, commitInfo)
	for parentHash := range commitInfo.HashParents.Iter() {
		changelog = completeTagCommitsForChangelog(changelog, parentHash, commitsMap)
	}
	return changelog
}

func writeChangelogs(changelogs map[*git.Tag][]*git.ConsolidatedCommit) error {
	for tagInfo, commits := range changelogs {
		err := writeChangelog(tagInfo, commits)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeChangelog(tagInfo *git.Tag, commits []*git.ConsolidatedCommit) error {
	tagname := cutTagNamePrefix(tagInfo.Name)
	filename := "CHANGELOG-" + tagname + ".md"
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not write CHANGELOG file %s\n - cause by: %w", filename, err)
	}
	defer outFile.Close()

	_, err = outFile.WriteString("# Changelog " + tagname + "\n")
	if err == nil {
		_, err = outFile.WriteString("\n")
		_, err = outFile.WriteString("Version " + tagname + "  \n")
		_, err = outFile.WriteString("Tag du *" + tagInfo.When.Format("02/01/2006 15:04:05") + "* par *" + tagInfo.Author + "*.  \n")
		_, err = outFile.WriteString(tagInfo.Subject + "  \n")
	}

	if err == nil {
		err = writeChangelogFeature(outFile, commits)
	}
	if err == nil {
		err = writeChangelogFix(outFile, commits)
	}
	if err == nil {
		err = writeChangelogOther(outFile, commits)
	}
	if err != nil {
		return fmt.Errorf("could not write CHANGELOG file %s\n - cause by: %w", err)
	}
	return nil
}

func writeChangelogFeature(out *os.File, commits []*git.ConsolidatedCommit) error {
	_, err := out.WriteString("## Features\n")
	if err == nil {
		for _, commit := range commits {
			if strings.HasPrefix(commit.Commit.Subject, "feat") {
				err = writeChangelogCommitLine(out, commit)
				if err != nil {
					return err
				}
			}
		}
	}
	if err == nil {
		_, err = out.WriteString("\n")
	}
	return err
}
func writeChangelogFix(out *os.File, commits []*git.ConsolidatedCommit) error {
	_, err := out.WriteString("## Fix\n")
	for _, commit := range commits {
		if strings.HasPrefix(commit.Commit.Subject, "fix") {
			err = writeChangelogCommitLine(out, commit)
			if err != nil {
				return err
			}
		}
	}
	if err == nil {
		_, err = out.WriteString("\n")
	}
	return err
}
func writeChangelogOther(out *os.File, commits []*git.ConsolidatedCommit) error {
	_, err := out.WriteString("## Other\n")
	for _, commit := range commits {
		if !strings.HasPrefix(commit.Commit.Subject, "feat") &&
			!strings.HasPrefix(commit.Commit.Subject, "fix") {
			err = writeChangelogCommitLine(out, commit)
			if err != nil {
				return err
			}
		}
	}
	if err == nil {
		_, err = out.WriteString("\n")
	}
	return err
}
func writeChangelogCommitLine(out *os.File, commit *git.ConsolidatedCommit) error {
	message := commit.Commit.Subject
	message = strings.ReplaceAll(message, "\n", "\n  ")
	message = replaceJiraNumberByJiraLink(message)
	_, err := out.WriteString("* " + message + "\n")
	return err
}
