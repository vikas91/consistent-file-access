package models

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Node struct {
	PeerId uuid.UUID `peerId`
	Address string `address`
	Balance float32 `balance`
	PublicKey rsa.PublicKey `publicKey`
}

type NodeList map[Node]uuid.UUID


func (nodeList *NodeList)GetNodeListJSON() string{
	nodeJSON, err := json.Marshal(nodeList)
	if(err!=nil){
		fmt.Println("Unable to convert application node list to json")
	}
	return string(nodeJSON)
}