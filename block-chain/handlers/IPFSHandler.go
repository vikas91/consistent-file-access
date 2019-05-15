package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/vikas91/consistent-file-access/block-chain/models"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const IPFS_DIR = "/tmp/ipfs"

// This will show the IPFS List available at peer node
func ShowIPFSList(w http.ResponseWriter, r *http.Request) {
	ipfsList = GetNodeIPFSList()
	ipfsListJSON := ipfsList.ShowNodeIPFSList()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ipfsListJSON))
}


// This will create a file to the IPFS directory
// This could be file added from application to block chain
func CreateIPFS(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("fileupload")
	defer file.Close()
	if err != nil {
		fmt.Println("Create file request experienced an error", err)
	}
	fileName := header.Filename
	fmt.Println("FileName", header.Filename)
	if _, err := os.Stat(path.Join(IPFS_DIR, fileName+"_version_1")); os.IsNotExist(err) {
		f, _ := os.OpenFile(path.Join(IPFS_DIR, fileName+"_version_1"), os.O_WRONLY|os.O_CREATE, 0444)
		defer f.Close()
		io.Copy(f, file)
	}else{
		fmt.Println("File already exists. Incrementing version and saving file")
		// This will check for all files with version pattern and then increment the version and save it
		pattern := path.Join(IPFS_DIR, fileName)+ "_version_*"
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Println("Error in matching file path", err)
		}
		nextVersionCount := strconv.Itoa(len(matches)+1)
		f, _ := os.OpenFile(path.Join(IPFS_DIR, fileName+"_version_"+ nextVersionCount), os.O_WRONLY|os.O_CREATE, 0444)
		defer f.Close()
		io.Copy(f, file)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
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
	if(signedIPFSHeartBeat.SignedNode.VerifyPeerSignature(signedIPFSHeartBeat.SignedIPFS, signedIPFSHeartBeat.IPFSListJSON)){
		ipfsList = GetNodeIPFSList()
		ipfsList.SyncNodeIPFSList(signedIPFSHeartBeat.IPFSListJSON)
		if(signedIPFSHeartBeat.Hops>0){
			fmt.Println("Forwarding ipfs heartbeats to all peers")
			peerList = GetPeerNodePeerList()
			peerNode = GetPeerNode()
			signedIPFSHeartBeat.Hops = signedIPFSHeartBeat.Hops - 1
			signedIPFSHeartBeat.ForwardNode = peerNode
			peerList.BroadcastSignedIPFSHeartBeat(signedIPFSHeartBeat)
		}
	}else{
		errorJSON = "Unable to sync ipfs heart beat data to node ipfs list"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(errorJSON))
}

// This will request IPFS File available at node which could be either ipfs file sharers or seeders
// If file is requested by not shared owners or seeders then return forbidden access request
func ShowIPFSFile(w http.ResponseWriter, r *http.Request) {

}

// This will request IPFS File available at node which could be either ipfs file sharers or seeders
func ShowIPFSSeedRequestList(w http.ResponseWriter, r *http.Request) {

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



