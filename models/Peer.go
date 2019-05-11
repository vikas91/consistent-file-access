package models

import (
	"crypto"
	"github.com/google/uuid"
	"sync"
)

type Peer struct {
	PeerId uuid.UUID
	Address string
	Balance float32
	PublicKey crypto.PublicKey
}

func NewPeer(address string , publicKey crypto.PublicKey) Peer{
	peerNode := Peer{ PeerId: uuid.New(), Address: address, PublicKey: publicKey, Balance: 0}
	return peerNode
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
	selfId uuid.UUID
	peerMap map[string]uuid.UUID
	maxLength int32
	mux sync.Mutex
}

// This will create a new peer list
func NewPeerList(id uuid.UUID, maxLength int32) PeerList {
	peerMap := make(map[string]uuid.UUID)
	peerList := PeerList{selfId: id, peerMap:peerMap, maxLength:maxLength}
	return peerList
}


