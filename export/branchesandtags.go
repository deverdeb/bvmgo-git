package export

import (
	"fmt"
	"github.com/deverdeb/bvmgo-git/git"
	"os"
)

func ExportBranchesAndTagsDiagram(commitsMap git.ConsolidatedCommitMap) error {
	cleanCommitMap := keepOnlyBranchesAnsTags(commitsMap)

	filename := "BRANCHES.md"
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not write file %s\n - cause by: %w", filename, err)
	}
	defer outFile.Close()

	outFile.WriteString("# GIT branches and tags\n")

	outFile.WriteString("```mermaid\n")
	outFile.WriteString("flowchart BT\n")
	for _, info := range cleanCommitMap {
		outFile.WriteString("    " + computeMermaidCommitMessage(info) + "\n")
	}
	for _, info := range cleanCommitMap {
		writeMermaidCommitRelation(info, outFile)
	}
	outFile.WriteString("```\n")

	return nil
}

func keepOnlyBranchesAnsTags(commitsMap git.ConsolidatedCommitMap) git.ConsolidatedCommitMap {
	cleanCommitMap := commitsMap.Clone()
	for hash, ConsolidatedCommit := range commitsMap {
		if !ConsolidatedCommit.IsTag() && !ConsolidatedCommit.IsBranch() {
			cleanCommitMap.RemoveByHash(hash)
		}
	}
	return cleanCommitMap
}
