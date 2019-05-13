package models

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Peer struct {
	PeerId uuid.UUID `peerId`
	Address string `address`
	Balance float32 `balance`
	PublicKey rsa.PublicKey `publicKey`
}

type SignedPeer struct{
	SignedPeerNode string
	PeerNode Peer
}

type PeerList map[uuid.UUID]Peer


func (node *Peer)GetNodeJSON() string{
	nodeJSON, err := json.Marshal(node)
	if(err!=nil){
		fmt.Println("Unable to convert application node to json", err)
	}
	return string(nodeJSON)
}

func (nodeList *PeerList)GetNodeListJSON() string{
	nodeListJSON, err := json.Marshal(nodeList)
	if(err!=nil){
		fmt.Println("Unable to convert application node list to json", err)
	}
	return string(nodeListJSON)
}


func (node *Peer)VerifyPeerSignature(signature string) bool{
	hexSignature, _ := hex.DecodeString(signature)
	nodeJSON := node.GetNodeJSON()
	byteNodeJSON := []byte(nodeJSON)
	hashed := sha256.Sum256(byteNodeJSON)
	err := rsa.VerifyPKCS1v15(&node.PublicKey, crypto.SHA256, hashed[:], hexSignature)
	if(err!=nil){
		return true
	}else{
		fmt.Println("signature verification failed", err)
	}
	return false
}
