package git

import (
	"fmt"
	"time"
)

type Hash string

func (hash Hash) String() string {
	return string(hash)
}

type Commit struct {
	Hash    Hash
	Parents []Hash
	When    time.Time
	Author  string
	Subject string
}

func (commit *Commit) String() string {
	return fmt.Sprintf("Commit[%s]{%s - %s}", commit.Hash, commit.Author, commit.When.Format("2006-01-02T15:04:05"))
}

type Branch struct {
	Name   string
	Commit *Commit
}

type Tag struct {
	Name    string
	When    time.Time
	Author  string
	Subject string
	Commit  *Commit
}

func (tag *Tag) String() string {
	if len(tag.Author) > 0 {
		return fmt.Sprintf("Tag[%s]{%s - %s}", tag.Name, tag.Author, tag.When.Format("2006-01-02T15:04:05"))
	} else {
		return fmt.Sprintf("Tag[%s]", tag.Name)
	}
}
