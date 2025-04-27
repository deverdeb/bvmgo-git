package export

import (
	"fmt"
	"github.com/deverdeb/bvmgo-git/git"
	"os"
	"strings"
)

func ExportMermaidGitTree(commitsMap map[git.Hash]*git.ConsolidatedCommit) error {
	filename := "Graph.md"
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not write file %s\n - cause by: %w", filename, err)
	}
	defer outFile.Close()

	outFile.WriteString("```mermaid\n")
	outFile.WriteString("flowchart BT\n")
	for _, info := range commitsMap {
		outFile.WriteString("    " + computeMermaidCommitMessage(info) + "\n")
	}
	for _, info := range commitsMap {
		writeMermaidCommitRelation(info, outFile)
	}
	outFile.WriteString("```\n")

	return nil
}

func computeMermaidCommitId(hash git.Hash) string {
	return "ID" + truncateHash(hash)
}

func computeMermaidCommitMessage(info *git.ConsolidatedCommit) string {
	begin := "(["
	end := "])"
	if info.IsBranch() {
		begin = "["
		end = "]"
	}
	message := computeMermaidCommitId(info.Commit.Hash) + begin + "\"`"
	for _, branchName := range info.Branches {
		message += "\n🌳**" + cutBranchNamePrefix(branchName) + "**"
	}
	for _, tagInfo := range info.Tags {
		message += "\n📂*" + cutTagNamePrefix(tagInfo.Name) + "*"
	}
	message += "\n" + truncateText(mermaidEscapeMessageCharacters(info.Commit.Subject), 40)
	message += "\n commit " + truncateHash(info.Commit.Hash)
	message += "\n" + info.Commit.When.Format("02/01/2006 15:04:05")
	message += "`\"" + end
	return message
}

func mermaidEscapeMessageCharacters(message string) string {
	result := strings.ReplaceAll(message, "`", "'")
	result = strings.ReplaceAll(result, "\"", "'")
	return result
}

func writeMermaidCommitRelation(info *git.ConsolidatedCommit, outFile *os.File) {
	idCommit := computeMermaidCommitId(info.Commit.Hash)
	for childHash := range info.HashChildren.Iter() {
		outFile.WriteString("    " + idCommit + " --> " + computeMermaidCommitId(childHash) + "\n")
	}
}
