package git

type ConsolidatedCommitMap map[Hash]*ConsolidatedCommit

func NewConsolidatedCommitMap() ConsolidatedCommitMap {
	return make(map[Hash]*ConsolidatedCommit)
}

func (commitMap ConsolidatedCommitMap) Clone() ConsolidatedCommitMap {
	clone := NewConsolidatedCommitMap()
	for hash, infoCommit := range commitMap {
		clone[hash] = infoCommit.Clone()
	}
	return clone
}

func (commitMap ConsolidatedCommitMap) FindByHash(hash Hash) *ConsolidatedCommit {
	info, _ := commitMap[hash]
	return info
}

func (commitMap ConsolidatedCommitMap) FindByHashes(hashes ...Hash) []*ConsolidatedCommit {
	infos := make([]*ConsolidatedCommit, 0)
	for _, hash := range hashes {
		info, _ := commitMap[hash]
		if info != nil {
			infos = append(infos, info)
		}
	}
	return infos
}

func (commitMap ConsolidatedCommitMap) RemoveByHash(hash Hash) {
	current := commitMap.FindByHash(hash)
	if current == nil {
		return
	}
	parents := commitMap.FindByHashes(current.HashParents.ToSlice()...)
	children := commitMap.FindByHashes(current.HashChildren.ToSlice()...)
	if len(parents) > 0 || len(children) > 0 {
		for _, parent := range parents {
			parent.HashChildren.Remove(hash)
			parent.HashChildren.AddAll(current.HashChildren.ToSlice()...)
		}
		for _, child := range children {
			child.HashParents.Remove(hash)
			child.HashParents.AddAll(current.HashParents.ToSlice()...)
		}
	}
	delete(commitMap, hash)
}
