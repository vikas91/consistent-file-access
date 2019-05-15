package models

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
)

type SignedIPFSHeartBeat struct {
	Node Peer `json:"peer"`
	SignedIPFS string `json:"signedIPFS"`
	IPFSListJSON string `json:"ipfsList"`
	Hops int32  `json:"hops"`
}

// This will create heart beat for ipfs created
func (peerNode *Peer) CreateSignedIPFSHeartBeat(peerNodeRSAKey *rsa.PrivateKey, ipfsListJSON string) SignedIPFSHeartBeat {
	signedSignture := peerNode.GetSignedSignature(peerNodeRSAKey, ipfsListJSON)
	signedIPFSHeartBeat := SignedIPFSHeartBeat{Node: *peerNode, SignedIPFS: signedSignture, IPFSListJSON: ipfsListJSON, Hops: 2}
	return signedIPFSHeartBeat
}

func (signedIPFSHeartBeat *SignedIPFSHeartBeat) GetSignedIPFSHeartBeatJSON() string{
	signedIPFSHeartBeatJSON, err := json.Marshal(signedIPFSHeartBeat)
	if(err!=nil){
		fmt.Println("Unable to convert signed heart beat to json")
	}
	return string(signedIPFSHeartBeatJSON)
}




