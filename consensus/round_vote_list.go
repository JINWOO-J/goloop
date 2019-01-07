package consensus

import (
	"github.com/icon-project/goloop/common"
)

type VoteItem struct {
	PrototypeIndex int16
	Signature      common.Signature
}

// TODO rename -> voteList
type roundVoteList struct {
	Prototypes []vote
	VoteItems  []VoteItem
}

func (vl *roundVoteList) AddVote(msg *voteMessage) {
	index := -1
	for i, p := range vl.Prototypes {
		if p.Equal(&msg.vote) {
			index = i
			break
		}
	}
	if index == -1 {
		vl.Prototypes = append(vl.Prototypes, msg.vote)
		index = len(vl.Prototypes) - 1
	}
	vl.VoteItems = append(vl.VoteItems, VoteItem{
		PrototypeIndex: int16(index),
		Signature:      msg.Signature,
	})
}

func (vl *roundVoteList) Len() int {
	return len(vl.VoteItems)
}

func (vl *roundVoteList) Get(i int) *voteMessage {
	msg := newVoteMessage()
	msg.vote = vl.Prototypes[vl.VoteItems[i].PrototypeIndex]
	msg.Signature = vl.VoteItems[i].Signature
	return msg
}

func newRoundVoteList() *roundVoteList {
	return &roundVoteList{}
}