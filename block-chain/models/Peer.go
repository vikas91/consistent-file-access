package models

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Peer struct {
	PeerId uuid.UUID `peerId`
	Address string `address`
	Balance float32 `balance`
	PublicKey rsa.PublicKey `publicKey`
	mux sync.Mutex
}

type SignedPeer struct{
	SignedPeerNode string
	PeerNode Peer
}

func NewPeer(address string , publicKey rsa.PublicKey) Peer{
	peerNode := Peer{ PeerId: uuid.New(), Address: address, PublicKey: publicKey, Balance: 0}
	return peerNode
}

func (peerNode *Peer)GetNodeJSON() string{
	peerNodeJSON, err := json.Marshal(peerNode)
	if(err!=nil){
		fmt.Println("Unable to convert peer node to json")
	}
	return string(peerNodeJSON)
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

func(peerNode *Peer) RegisterPeer(registerURL string) map[string]uuid.UUID {
	peerNode.mux.Lock()

	peerNode.mux.Unlock()
}

// This will update peerNode peerMap with peerMap from register server
func(peerList *PeerList) UpdatePeerList(peerMap map[string]uuid.UUID){
	peerList.mux.Lock()
	for key, _ := range peerMap {
		value, ok := peerList.peerMap[key]
		if !ok {
			peerList.peerMap[key] = value
		}
	}
	peerList.mux.Unlock()
}
