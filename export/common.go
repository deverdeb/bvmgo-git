package export

import (
	"fmt"
	"github.com/deverdeb/bvmgo-git/git"
	"github.com/deverdeb/bvmgo-git/utils/sets"
	"maps"
	"regexp"
	"strings"
)

var jiraIdRegex = regexp.MustCompile("(AG-[0-9]{4})")
var jiraBaseUrl = "https://astree-software.atlassian.net/browse/"

func replaceJiraNumberByJiraLink(message string) string {
	jiraIds := jiraIdRegex.FindAllString(message, -1)
	jiraIdSet := sets.NewFromSlice[string](jiraIds)
	newMessage := message
	for jiraId := range jiraIdSet.Iter() {
		jiraUrl := jiraBaseUrl + jiraId
		markdownLink := fmt.Sprintf("[%s](%s)", jiraId, jiraUrl)
		newMessage = strings.ReplaceAll(newMessage, jiraId, markdownLink)
	}
	return newMessage
}

func truncateText(text string, max int) string {
	currentText := text
	newLineIndex := strings.Index(text, "\n")
	if newLineIndex != -1 {
		currentText = text[:newLineIndex]
	}
	if max > len(currentText) {
		return currentText
	}
	return currentText[:strings.LastIndex(currentText[:max-3], " ")] + "..."
}

func truncateHash(hash git.Hash) string {
	hashText := hash.String()
	if len(hashText) <= 7 {
		return hashText
	}
	return hashText[:7]
}

func cutBranchNamePrefix(initialBranchName string) string {
	branchName, _ := strings.CutPrefix(initialBranchName, "refs/heads/")
	branchName, _ = strings.CutPrefix(branchName, "refs/remotes/")
	return branchName
}

func cutTagNamePrefix(initialTagName string) string {
	tagName, _ := strings.CutPrefix(initialTagName, "refs/tags/")
	return tagName
}

func GroupRemovableCommit(commitList git.ConsolidatedCommitMap) git.ConsolidatedCommitMap {
	reducedCommitList := commitList.Clone()
	for hash := range maps.Keys(reducedCommitList) {
		info := reducedCommitList[hash]
		if info == nil ||
			!isRemovableCommit(info) {
			continue
		}
		removeCommit(reducedCommitList, info)
	}
	return reducedCommitList
}

func removeCommit(commitList git.ConsolidatedCommitMap, commitInfo *git.ConsolidatedCommit) {
	current := commitInfo
	parent := commitList[current.HashParents.ToSlice()[0]]
	child := commitList[current.HashChildren.ToSlice()[0]]
	if !isRemovableCommit(parent) && !isRemovableCommit(child) {
		// Keep commit
		return
	}

	for isRemovableCommit(parent) {
		current = parent
		parent = commitList[current.HashParents.ToSlice()[0]]
	}
	child = commitList[current.HashChildren.ToSlice()[0]]
	nbCommits := 1
	for isRemovableCommit(child) {
		// Virer le commit de la map
		maps.DeleteFunc(commitList, func(hash git.Hash, infoCheckDelete *git.ConsolidatedCommit) bool {
			return infoCheckDelete == current
		})
		// Refaire les liaisons parent / enfant sans le commit
		if parent != nil {
			parent.HashChildren.Remove(current.Commit.Hash)
			if child != nil {
				parent.HashChildren.Add(child.Commit.Hash)
			}
		}
		if child != nil {
			child.HashParents.Remove(current.Commit.Hash)
			if parent != nil {
				child.HashParents.Add(parent.Commit.Hash)
			}
		}
		// Passer à la suite
		nbCommits++
		current = child
		child = commitList[current.HashChildren.ToSlice()[0]]
	}
	// Transformer le commit courant en commit de "regroupement"
	current.Commit.Subject = fmt.Sprintf("(%d commits)", nbCommits)
}

func isRemovableCommit(commit *git.ConsolidatedCommit) bool {
	if commit == nil {
		return false
	}
	// commit supprimable = sans branche, sans tag et avec 1 parent / 1 enfant
	hasOneChild := commit.HashChildren.Size() == 1
	hasOneParent := commit.HashParents.Size() == 1
	hasNotTagAndBranch := !commit.IsTag() && !commit.IsBranch()
	return hasOneChild && hasOneParent && hasNotTagAndBranch
}
