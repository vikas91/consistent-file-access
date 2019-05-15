package handlers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vikas91/consistent-file-access/block-chain/models"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

// This should come from config parameters
const REGISTER_ADDR = "http://localhost:6686"
const REGISTER_URL = REGISTER_ADDR + "/register"
var peerNodeRSAKey *rsa.PrivateKey

var peerNode models.Peer
var peerList models.PeerList
var blockChain models.BlockChain
var ipfsList models.IPFSList
var ifStarted bool

func GetNodeIPFSList() models.IPFSList {
	return ipfsList
}

func GetPeerNode() models.Peer {
	return peerNode
}

func GetPeerNodePeerList() models.PeerList {
	return peerList
}

func InitializePeerNode(args []string){
	var nodePort int32
	if len(os.Args) > 1 {
		i, err := strconv.ParseInt(os.Args[1], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		nodePort = int32(i)

	} else {
		nodePort = 8000
	}
	newPeerList := Register(nodePort)
	peerList = models.NewPeerList(peerNode.PeerId, 32)
	peerList.UpdatePeerList(newPeerList)
	blockChain = models.NewBlockChain()
	ipfsList = models.NewIPFSList()
	ipfsList.FetchNodeIPFSList(peerNode)
}

// This will return the IP Address of Node
func generateNodeIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// This will generate an RSA key pair for the node
func generateNodeKeyPair() *rsa.PrivateKey {
	reader := rand.Reader
	privateKey, err := rsa.GenerateKey(reader, 4096)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	peerNodeRSAKey = privateKey
	return privateKey
}


// This will create key pair for node and create the peer node
// This will also register the public key of node on application
func Register(port int32) map[uuid.UUID]models.Peer {
	ipAddress := generateNodeIPAddress()
	rsaPrivateKey := generateNodeKeyPair()
	publicKey := rsaPrivateKey.PublicKey
 	completeAddress := ipAddress + ":" + fmt.Sprint(port)
	peerNode = models.NewPeer(completeAddress, publicKey)
	newPeerList := peerNode.RegisterPeer(peerNodeRSAKey, REGISTER_URL)
	return newPeerList
}

// This will create a new key pair for node
// This will also register the new public key of node on application
// This will send a signed PeerNode with old private key to application to update
func UpdatePeerNodeKeyPair(){
	reader := rand.Reader
	privateKey, err := rsa.GenerateKey(reader, 4096)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	oldRSAKey := peerNodeRSAKey
	tempPeerNode := peerNode
	tempPeerNode.PublicKey = privateKey.PublicKey
	peerJSON := tempPeerNode.GetNodeJSON()
	hashed := sha256.Sum256([]byte(peerJSON))
	signature, err := rsa.SignPKCS1v15(reader, oldRSAKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}

	fmt.Printf("Signature: %x\n", signature)
	//TODO: Call Application Register Peer with Peer Node Json and signed signature

	peerNodeRSAKey = privateKey
	peerNode.PublicKey = privateKey.PublicKey
}

func GetIPFSListJSON(ipfsList []models.IPFS) string{
	ipfsListJSON, err := json.Marshal(ipfsList)
	if(err!=nil){
		fmt.Println("Unable to convert new ipfs list node to json")
	}
	return string(ipfsListJSON)
}

// This will periodically check for new files and update the IPFS list in directory
func PeriodicUpdateNodeIPFSList(){
	for ifStarted {
		newIPFSList := ipfsList.PollNodeIPFSList(peerNode)
		newIPFSListJSON := GetIPFSListJSON(newIPFSList)
		signedIPFSHeartBeat := peerNode.CreateSignedIPFSHeartBeat(peerNodeRSAKey, newIPFSListJSON)
		peerList.BroadcastSignedIPFSHeartBeat(signedIPFSHeartBeat)
	}
}

// This will start the node server
// It will first update the block chain from register address host
// This will run a thread to get the file in IPFS directory periodically and send IPFS heartbeat to peers
func StartNode(w http.ResponseWriter, r *http.Request) {
	if (!ifStarted){
		fmt.Println("Starting Peer Node")
		ifStarted = true
		go func() {
			PeriodicUpdateNodeIPFSList()
		}()
		fmt.Println("Started Peer Node")
	}
	GetNodeDetails(w,r)
}

// This will stop the peer node from generating new data
func StopNode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stopping Peer Node")
	ifStarted = false
	fmt.Println("Stopped Peer Node")
	GetNodeDetails(w, r)
}


// This will get the details of peer Node
func GetNodeDetails(w http.ResponseWriter, r *http.Request) {
	peerJSON := peerNode.GetNodeJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peerJSON))
}

// This will get the details of peer Node
func GetNodePeerList(w http.ResponseWriter, r *http.Request) {
	peerListJSON := peerList.GetPeerListJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peerListJSON))
}

// This will update the details of peer Node
func UpdateNodeDetails(w http.ResponseWriter, r *http.Request) {
	// TODO: This should not be publicly accesible. Need to think of better way to handle this
	UpdatePeerNodeKeyPair()
	GetNodeDetails(w, r)
}
