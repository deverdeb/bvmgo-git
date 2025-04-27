package main

import (
	"fmt"
	"github.com/deverdeb/bvmgo-git/export"
	"github.com/deverdeb/bvmgo-git/git"
	"github.com/deverdeb/bvmgo-term/ansi"
	"github.com/deverdeb/bvmgo-term/term"
	"log"
)

func main() {
	blueTitle := term.Style{Foreground: &term.Blue, Uppercase: true}
	ansi.ClearScreen()
	fmt.Println(blueTitle.Sprintf("\nGit tree extraction"))

	commitsMap, err := git.ExtractCommits()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(blueTitle.Sprintf("\nGenerate changelogs"))

	err = export.ExportChangelogs(commitsMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(blueTitle.Sprintf("\nGenerate branches and tags diagram"))
	err = export.ExportBranchesAndTagsDiagram(commitsMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(blueTitle.Sprintf("\nGenerate all tree diagram"))
	reducedCommitsMap := export.GroupRemovableCommit(commitsMap)
	err = export.ExportMermaidGitTree(reducedCommitsMap)
	if err != nil {
		log.Fatal(err)
	}

}
