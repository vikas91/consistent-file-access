package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vikas91/consistent-file-access/application/models"
	"net/http"
	"path"
)


const TEMPLATE_DIR = "/Users/victor/workspace/go/src/github.com/vikas91/consistent-file-access/application/templates/"

var nodeList models.PeerList

func InitializeApplication(){
	nodeList = make(map[uuid.UUID]models.Peer)
}

func Index(w http.ResponseWriter, r *http.Request) {
	p := path.Dir(TEMPLATE_DIR+"index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

// This will add the peer node to register server node list
// Returns the node list json as response
func RegisterNode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register Node called")
	decoder := json.NewDecoder(r.Body)
	var signedNode models.SignedPeer
	err := decoder.Decode(&signedNode)
	if err != nil {
		fmt.Println("Unable to decode request to register node. Incorrect format")
	}
	var nodeListJSON string
	peerNode, ok := nodeList[signedNode.PeerNode.PeerId]
	if !ok {
		if(signedNode.PeerNode.VerifyPeerSignature(signedNode.SignedPeerNode)){
			fmt.Println("Signature Verified for new user. Register Successful")
			nodeList[signedNode.PeerNode.PeerId] = signedNode.PeerNode
			nodeListJSON = nodeList.GetNodeListJSON()
		}else{
			nodeListJSON = "{}"
		}
	}else{
		fmt.Println("Peer Id already registered with this uuid.")
		if(peerNode.VerifyPeerSignature(signedNode.SignedPeerNode)){
			fmt.Println("Verified with old public key. Updating key pair")
			nodeList[signedNode.PeerNode.PeerId] = signedNode.PeerNode
			nodeListJSON = nodeList.GetNodeListJSON()
		}else{
			nodeListJSON = "{}"
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nodeListJSON))
}


// This will show list of all nodes registered on the server
func ShowNodeList(w http.ResponseWriter, r *http.Request) {
	peerJSON := nodeList.GetNodeListJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peerJSON))
}
