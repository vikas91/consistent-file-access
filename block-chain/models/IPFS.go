package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"
)

const IPFS_DIR = "/tmp/ipfs"

type IPFS struct {
	Id uuid.UUID
	FileName string
	Address string
	FileVersionList []IPFSVersion
}

type IPFSVersion struct {
	Id int32
	PreviousVersionHash string
	CurrentVersionHash string
	SeedCost int32
	SeedCount int32
	SeedEnabled bool
	VersionOwners []IPFSUser
	VersionSeeders []IPFSUser
	CreatedTime time.Time
}

type IPFSUser struct {
	PeerNode Peer
	PeerFileKey string
}

type IPFSList struct {
	IPFSMap map[uuid.UUID]IPFS
	UpdatedTime time.Time
	mux sync.Mutex
}

func RandomSeed(){
	rand.Seed(time.Now().UTC().UnixNano())
}

// Sleep for random time between 5 to 10seconds
func RandomSleep(){
	RandomSeed()
	sleepTime := 5 + rand.Intn(5)
	fmt.Println("Sleeping for:", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)
}

// This will create a new ipfs entry
func NewIPFS(file os.FileInfo) IPFS {
	fileName := file.Name()
	absoluteFilePath := path.Join(IPFS_DIR, fileName)
	fileStat, _ := os.Stat(absoluteFilePath)
	mtime := fileStat.ModTime()

	fileHash, err := FileMD5Hash(absoluteFilePath)
	if(err!=nil){

	}
	ipfsVersion := IPFSVersion{Id: 1, PreviousVersionHash: "root", CurrentVersionHash: fileHash, CreatedTime: mtime}
	ipfs := IPFS{Id: uuid.New(), FileName: file.Name(), FileVersionList: []IPFSVersion{ipfsVersion}}
	return ipfs
}

// This will create a new ipfs list
func NewIPFSList() IPFSList {
	ipfsMap := make(map[uuid.UUID]IPFS)
	ipfsList := IPFSList{IPFSMap: ipfsMap, UpdatedTime: time.Now()}
	return ipfsList
}


// This will get MD5Hash of file given absolute file path
func FileMD5Hash(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}

// This is called initially when node is started
// This will update all files present in IPFS directory to Peer IPFS List
func (ipfsList *IPFSList)FetchNodeIPFSList(){
	ipfsList.mux.Lock()
	defer ipfsList.mux.Unlock()

	// This will create IPFS directory if it does not exist
	if _, err := os.Stat(IPFS_DIR); os.IsNotExist(err) {
		os.Mkdir(IPFS_DIR, 0700)
	}

	files, _ := ioutil.ReadDir(IPFS_DIR)
	for _, file := range files {
		ipfs := NewIPFS(file)
		ipfsList.IPFSMap[ipfs.Id] = ipfs
	}
}

// This will update the node ipfs list periodically
// Checks for newly created /modified files and then creates a new entry to IPFS list
func (ipfsList *IPFSList)PollNodeIPFSList(){
	ipfsList.mux.Lock()
	defer ipfsList.mux.Unlock()
	RandomSleep()
	files, _ := ioutil.ReadDir(IPFS_DIR)
	for _, file := range files {
		fileName := file.Name()
		absoluteFilePath := path.Join(IPFS_DIR, fileName)
		fileStat, _ := os.Stat(absoluteFilePath)
		mtime := fileStat.ModTime()
		if(mtime.After(ipfsList.UpdatedTime)){
			//TODO: Take care of new file and updated file logic here
			// TODO: Should send new IPFS HEARTBEAT with signed signature for both scenarios
			ipfs := NewIPFS(file)
			fmt.Println("New IPFS File Found", ipfs.FileName, ipfs.FileVersionList, "at", ipfsList.UpdatedTime)
			ipfsList.IPFSMap[ipfs.Id] = ipfs
		}
	}
	fmt.Println("All files have been scanned at", time.Now())
	ipfsList.UpdatedTime = time.Now()
}


func (ipfsList *IPFSList)GetNodeIPFSJSON() string{
	ipfsListJSON, err := json.Marshal(ipfsList)
	if(err!=nil){
		fmt.Println("Unable to convert peer node ipfs list to json")
	}
	return string(ipfsListJSON)
}

// This will return the IPFS List of peer as json
func (ipfsList *IPFSList)ShowNodeIPFSList() string{
	ipfsList.mux.Lock()
	defer ipfsList.mux.Unlock()
	return ipfsList.GetNodeIPFSJSON()
}