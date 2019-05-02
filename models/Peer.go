package models

import (
	"sync"
)

type Peer struct {
	PeerId int32
	Address string
	Balance float32
	PublicKey string
}

type PeerTransactionList struct {
	PeerNode Peer
	TransactionMap map[string]Transaction
}

type PeerIPFSList struct{
	PeerNode Peer
	IPFSMap map[int32]IPFS
}

type PeerIPFSPendingShareList struct{
	IPFSMap map[string]IPFS
}

type PeerList struct {
	selfId int32
	peerMap map[string]int32
	maxLength int32
	mux sync.Mutex
}

// This will create a new peer list
func NewPeerList(id int32, maxLength int32) PeerList {
	peerMap := make(map[string]int32)
	peerList := PeerList{selfId: id, peerMap:peerMap, maxLength:maxLength}
	return peerList
}


