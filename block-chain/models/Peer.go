package models

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
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
	peerNode.mux.Lock()
	defer peerNode.mux.Unlock()
	peerNodeJSON, err := json.Marshal(peerNode)
	if(err!=nil){
		fmt.Println("Unable to convert peer node to json")
	}
	return string(peerNodeJSON)
}

func (peerNode *Peer) GetSignedSignature(priv *rsa.PrivateKey, message string) string{
	rng := rand.Reader
	byteMessage := []byte(message)
	hashed := sha256.Sum256(byteMessage)
	signature, err := rsa.SignPKCS1v15(rng, priv, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
	}
	sig := base64.StdEncoding.EncodeToString(signature)
	return sig
}

// This will register peer on register server. It returns a list of peer nodes
func(peerNode *Peer) RegisterPeer(priv *rsa.PrivateKey , registerURL string) map[uuid.UUID]Peer {
	peerNode.mux.Lock()
	defer peerNode.mux.Unlock()

	// TODO: Send signed signature for registration to prevent spam

	peerJSON, err := json.Marshal(peerNode)
	signedSignature := peerNode.GetSignedSignature(priv, string(peerJSON))
	signedPeer := SignedPeer{SignedPeerNode: signedSignature, PeerNode: *peerNode}

	signedPeerJSON, err := json.Marshal(signedPeer)

	fmt.Println("Initiating connection to Register Server")
	if err != nil {
		fmt.Println("Unable to convert peer node to json")
	}
	response, err := http.Post(registerURL, "application/json", bytes.NewBuffer(signedPeerJSON))
	if(err!=nil){
		fmt.Println("Unable to connect to Register server. Service unavailable")
		return make(map[uuid.UUID]Peer)
	}else{
		fmt.Println("Connected to Register Server")
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		peerMapJSON := buf.String()
		peerMap := GetPeerListFromJSON(peerMapJSON)
		return peerMap
	}
}

type PeerList struct {
	selfId uuid.UUID
	peerMap map[uuid.UUID]Peer
	maxLength int32
	mux sync.Mutex
}

// This will create a new peer list
func NewPeerList(id uuid.UUID, maxLength int32) PeerList {
	peerMap := make(map[uuid.UUID]Peer)
	peerList := PeerList{selfId: id, peerMap:peerMap, maxLength:maxLength}
	return peerList
}
// This will update peerNode peerMap with peerMap from register server
func(peerList *PeerList) GetPeerListJSON() string{
	peerList.mux.Lock()
	defer peerList.mux.Unlock()
	peerListJSON, err := json.Marshal(peerList)
	if(err!=nil){
		fmt.Println("Unable to convert peer node to json")
	}
	return string(peerListJSON)
}

// This will convert JSON String to peermap
func GetPeerListFromJSON(peerListJSON string)map[uuid.UUID]Peer{
	var newPeerMap map[uuid.UUID]Peer
	err := json.Unmarshal([]byte(peerListJSON), &newPeerMap)
	if (err != nil) {
		fmt.Println(err)
	}
	return newPeerMap
}

// This will update peerNode peerMap with peerMap from register server
func(peerList *PeerList) UpdatePeerList(peerMap map[uuid.UUID]Peer){
	peerList.mux.Lock()
	for key, _ := range peerMap {
		if key != peerList.selfId {
			value, ok := peerList.peerMap[key]
			if !ok {
				peerList.peerMap[key] = value
			}
		}
	}
	peerList.mux.Unlock()
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
