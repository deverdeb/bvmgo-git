package git

import (
	"github.com/deverdeb/bvmgo-git/utils/sets"
)

type ConsolidatedCommit struct {
	Commit       *Commit
	HashParents  sets.Set[Hash]
	HashChildren sets.Set[Hash]
	Branches     []string
	Tags         []*Tag
}

func BuildConsolidatedCommit(commit *Commit) *ConsolidatedCommit {
	return &ConsolidatedCommit{
		Commit:       commit,
		HashParents:  sets.New[Hash](),
		HashChildren: sets.New[Hash](),
		Branches:     make([]string, 0),
		Tags:         make([]*Tag, 0),
	}
}

func (infoCommit *ConsolidatedCommit) Clone() *ConsolidatedCommit {
	clone := &ConsolidatedCommit{
		Commit:       infoCommit.Commit,
		HashParents:  infoCommit.HashParents.Clone(),
		HashChildren: infoCommit.HashChildren.Clone(),
		Branches:     make([]string, len(infoCommit.Branches)),
		Tags:         make([]*Tag, len(infoCommit.Tags)),
	}
	copy(clone.Branches, infoCommit.Branches)
	copy(clone.Tags, infoCommit.Tags)
	return clone
}

func (infoCommit *ConsolidatedCommit) IsBranch() bool {
	if infoCommit == nil {
		return false
	}
	return len(infoCommit.Branches) > 0
}

func (infoCommit *ConsolidatedCommit) IsTag() bool {
	if infoCommit == nil {
		return false
	}
	return len(infoCommit.Tags) > 0
}
