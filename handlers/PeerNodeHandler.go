package handlers

import (
	"fmt"
	"github.com/vikas91/consistent-file-access/models"
	"net"
	"net/http"
	"os"
	"strconv"
)

const REGISTER_ADDR = "http://localhost:6686"
var SELF_ADDR = "http://localhost:6686"

var peerList models.PeerList
var blockChain models.BlockChain
var ipfsList models.IPFSList
var ifStarted bool

func getNodeIPFSList() models.IPFSList {
	return ipfsList
}

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
	var selfId int32
	if len(os.Args) > 1 {
		i, err := strconv.ParseInt(os.Args[1], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		selfId = int32(i)

	} else {
		selfId = 6686
	}

	peerList = models.NewPeerList(selfId, 32)
	blockChain = models.NewBlockChain()
	ipfsList = models.NewIPFSList()
	ipfsList.FetchPeerNodeIPFSList()
}

// This will periodically check for new files and update the IPFS list in directory
func PeriodicUpdatePeerNodeIPFSList(){
	for ifStarted {
		ipfsList.PollPeerNodeIPFSList()
	}
}

// This will start the node server
// It will first update the block chain from register address host
// This will run a thread to get the file in IPFS directory periodically and send IPFS heartbeat to peers
func StartNode(w http.ResponseWriter, r *http.Request) {
	if (!ifStarted){
		fmt.Println("Starting Peer Node")
		ip, port, err := net.SplitHostPort(r.Host)
		userIP := net.ParseIP(ip)
		fmt.Println("Ip, Port, err", ip, port, err, userIP)
		ifStarted = true
		go func() {
			PeriodicUpdatePeerNodeIPFSList()
		}()
		fmt.Println("Started Peer Node")
	}
}

// This will stop the peer node from generating new data
func StopNode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stopping Peer Node")
	ifStarted = false
	fmt.Println("Stopped Peer Node")
}

// This will restart the peer node
func RestartNode(w http.ResponseWriter, r *http.Request) {
	StopNode(w , r)
	StartNode(w, r)
}
