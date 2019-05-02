package handlers

import (
	"fmt"
	"github.com/vikas91/consistent-file-access/models"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const REGISTER_ADDR = "http://localhost:6686"
var SELF_ADDR = "http://localhost:6686"

var peerList models.PeerList
var blockChain models.BlockChain
var ipfsList models.IPFSList
var ifStarted bool

func RandomSeed(){
	rand.Seed(time.Now().UTC().UnixNano())
}

// Sleep for random time between 5 to 10seconds
func RandomSleep(){
	RandomSeed()
	sleepTime := 5 + rand.Intn(5)
	time.Sleep(time.Duration(sleepTime) * time.Second)
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
}

// This will periodically check for new files and update the IPFS list in directory
func UpdateNodeIPFSList(ipfsList models.IPFSList){
	RandomSleep()
	models.UpdateNodeIPFSList(ipfsList)
}

// This will start the node server
// It will first update the block chain from register address host
// This will run a thread to get the file in IPFS directory periodically and send IPFS heartbeat to peers
func StartNode(w http.ResponseWriter, r *http.Request) {
	if (!ifStarted){
		ifStarted = true
		SELF_ADDR = r.Host
		fmt.Println("Host Address", SELF_ADDR)
		// Download complete block chain from remote host
		if(SELF_ADDR!=REGISTER_ADDR) {

		}
		for ifStarted {
			UpdateNodeIPFSList(ipfsList)
		}
	}
}

// This will stop the peer node from generating new data
func StopNode(w http.ResponseWriter, r *http.Request) {
	ifStarted = false
}

// This will restart the peer node
func RestartNode(w http.ResponseWriter, r *http.Request) {

}
