package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
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
	FileVersion []IPFSVersion
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

// This will create a new ipfs entry
func NewIPFS(file os.FileInfo) IPFS {
	fileName := file.Name()
	fileStat, _ := os.Stat(path.Join(IPFS_DIR, fileName))
	mtime := fileStat.ModTime()

	fileHash, _ := FileMD5Hash(path.Join(IPFS_DIR, file.Name()))
	ipfsVersion := IPFSVersion{Id: 1, PreviousVersionHash: "root", CurrentVersionHash: fileHash, CreatedTime: mtime}

	ipfs := IPFS{Id: uuid.New(), FileName: file.Name(), FileVersion: []IPFSVersion{ipfsVersion}}
	return ipfs
}

// This will create a new ipfs list
func NewIPFSList() IPFSList {
	ipfsMap := make(map[uuid.UUID]IPFS)
	ipfsList := IPFSList{IPFSMap: ipfsMap, UpdatedTime: time.Now()}
	return ipfsList
}


// This will get MD5Hash of file
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

// This will update the node ipfs list
// Checks for newly created /modified files and then creates a new entry to IPFS list
func UpdateNodeIPFSList(ipfsList IPFSList){
	if _, err := os.Stat(IPFS_DIR); os.IsNotExist(err) {
		os.Mkdir(IPFS_DIR, 0700)
	}
	files, _ := ioutil.ReadDir(IPFS_DIR)
	for _, file := range files {
		fileName := file.Name()
		fileStat, _ := os.Stat(path.Join(IPFS_DIR, fileName))
		mtime := fileStat.ModTime()
		if(mtime.After(ipfsList.UpdatedTime)){
			//TODO: Take care of new file and updated file logic here
			ipfs := NewIPFS(file)
			ipfsList.UpdatedTime = mtime
			ipfsList.IPFSMap[ipfs.Id] = ipfs
		}
	}
	fmt.Println(ipfsList)
}
