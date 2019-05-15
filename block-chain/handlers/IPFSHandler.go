package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/vikas91/consistent-file-access/block-chain/models"
	"net/http"
)


// This will show the IPFS List available at peer node
func ShowIPFSList(w http.ResponseWriter, r *http.Request) {
	ipfsList = GetNodeIPFSList()
	ipfsListJSON := ipfsList.ShowNodeIPFSList()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ipfsListJSON))
}


// This will sync the ipfs list received from peer node to its own ipfs list
// This will broadcast IPFSHeartBeat data to nearest peers available to a node
func IPFSHeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Heart Beat Receive Node called")
	decoder := json.NewDecoder(r.Body)
	var signedIPFSHeartBeat models.SignedIPFSHeartBeat
	err := decoder.Decode(&signedIPFSHeartBeat)
	if err != nil {
		fmt.Println("Unable to decode signed heart beat request. Incorrect format", err)
	}
	var errorJSON string
	if(signedIPFSHeartBeat.Node.VerifyPeerSignature(signedIPFSHeartBeat.SignedIPFS, signedIPFSHeartBeat.IPFSListJSON)){
		ipfsList = GetNodeIPFSList()
		ipfsList.SyncNodeIPFSList(signedIPFSHeartBeat.IPFSListJSON)
		if(signedIPFSHeartBeat.Hops>0){
			fmt.Println("Forwarding ipfs heartbeats to all peers")
			peerList = GetPeerNodePeerList()
			signedIPFSHeartBeat.Hops = signedIPFSHeartBeat.Hops - 1
			peerList.BroadcastSignedIPFSHeartBeat(signedIPFSHeartBeat)
		}
	}else{
		errorJSON := "Unable to sync ipfs heart beat data to node ipfs list"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(errorJSON))
}

// This will request IPFS File available at node which could be either ipfs file sharers or seeders
func ShowIPFSFile(w http.ResponseWriter, r *http.Request) {

}

// This will get the file versions available for an ipfs file
func ShowIPFSFileVersions(w http.ResponseWriter, r *http.Request) {

}

// This will show the current ipfs file who request for seeds in the networks
func ShowIPFSSeedRequests(w http.ResponseWriter, r *http.Request) {

}

// This is where the share request to file will be send to for every node
func ShareRequestIPFSFile(w http.ResponseWriter, r *http.Request) {

}

// This is where the seed accept request to file will be send to for seeds request send in IPFS heart bear
func SeedRequestIPFSFile(w http.ResponseWriter, r *http.Request) {

}



