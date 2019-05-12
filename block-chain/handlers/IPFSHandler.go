package handlers

import "net/http"


// This will show the IPFS List available at peer node
func ShowIPFSList(w http.ResponseWriter, r *http.Request) {
	ipfsList = getNodeIPFSList()
	ipfsListJSON := ipfsList.ShowNodeIPFSList()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ipfsListJSON))
}

// This will download and update the IPFS List from peer node
func UpdateIPFSList(w http.ResponseWriter, r *http.Request) {

}

// This will broadcast IPFSHeartBeat data to nearest peers available to a node
func IPFSHeartBeatReceive(w http.ResponseWriter, r *http.Request) {

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



