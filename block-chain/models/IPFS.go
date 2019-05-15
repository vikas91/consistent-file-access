package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	random "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IPFS struct {
	Id uuid.UUID
	FileName string
	FileVersionList []IPFSVersion
	mux sync.Mutex
}

type IPFSVersion struct {
	Id int
	VersionHash string
	SeedCost float32
	SeedCount int
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

// This will encrypt the AES File Key with peer node public key
func EncryptAESKey(peerNode Peer, aesFileKey string) string {
	byteArray := []byte(aesFileKey)
	label := []byte("aes encrypt file key")
	encryptedText, err := rsa.EncryptOAEP(sha256.New(), random.Reader, &peerNode.PublicKey , byteArray, label)
	if err != nil {
		fmt.Println("Unable to encrpyt AES Key of file", aesFileKey)
	}
	return string(encryptedText)
}

// This will encrypt the AES File Key with peer node public key
func DecryptAESKey(priv *rsa.PrivateKey, aesEncryptFileKey string) string {
	byteArray := []byte(aesEncryptFileKey)
	label := []byte("aes decrypt file key")
	decryptedText, err := rsa.DecryptOAEP(sha256.New(), random.Reader, priv , byteArray, label)
	if err != nil {
		fmt.Println("Unable to encrpyt AES Key of file", aesEncryptFileKey)
	}
	return string(decryptedText)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(random.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

// This will get MD5Hash of file given absolute file path
func FileMD5Hash(aesPassword string, filePath string) string {
	data := AESEncryptIPFSFileContent(aesPassword, filePath)
	var returnMD5String string
	hash := md5.New()
	io.WriteString(hash, data)
	hashInBytes := hash.Sum(nil)
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String
}

// This will encrypt contents of file using AES encryption
// Return AES passphrase for the file
func AESEncryptIPFSFileContent(aesPassword string, absoluteFilePath string) string {
	data, err := ioutil.ReadFile(absoluteFilePath)
	if(err!=nil){
		fmt.Println("Unable to encrypt file content", err)
	}
	ciphertext := encrypt(data, aesPassword)
	return string(ciphertext)
}


// This will decrypt contents of file given the AES passphrase key
func AESDecryptIPFSFile(absoluteFilePath string, aesPassword string){
	data, _ := ioutil.ReadFile(absoluteFilePath)
	filetext := decrypt(data, aesPassword)
	err := ioutil.WriteFile(absoluteFilePath, filetext, 0644)
	if(err!=nil){
		fmt.Println("unable to write to file", absoluteFilePath, err)
	}
}


func versionFileNameParser(fileName string)(string, int){
	s := strings.Split(fileName, "_version_")
	nonVersionedFileName  := s[0]
	versionNumber, _ := strconv.Atoi(s[1])
	return nonVersionedFileName, versionNumber
}

// This will create a new ipfs entry
// This will take filePath and peerNode as parameters
// AES encrypt file content and then get hash of encrypted file content
// The hash of encrypted file content will be used by miners to verify authenticity of file data
func NewIPFS(IPFS_DIR string, file os.FileInfo, peerNode Peer, ipfsList *IPFSList) IPFS {
	fileName := file.Name()
	absoluteFilePath := path.Join(IPFS_DIR, fileName)
	fileStat, _ := os.Stat(absoluteFilePath)
	mtime := fileStat.ModTime()
	aesPassphrase := uuid.New().String()
	encryptedAESKey := EncryptAESKey(peerNode, aesPassphrase)
	fileHash := FileMD5Hash(aesPassphrase, absoluteFilePath)
	ipfsUser := IPFSUser{PeerNode: peerNode, PeerFileKey: encryptedAESKey}

	nonVersionedFileName, versionNumber := versionFileNameParser(fileName)
	fmt.Println("FileName without version", nonVersionedFileName, versionNumber)
	ipfsVersion := IPFSVersion{Id: versionNumber, VersionHash: fileHash, VersionOwners: []IPFSUser{ipfsUser}, CreatedTime: mtime}
	prevVersionExists := false
	var newIPFS IPFS
	for _, ipfs := range ipfsList.IPFSMap {
		if(ipfs.FileName == nonVersionedFileName){
			prevVersionExists = true
			ipfs.FileVersionList = append(ipfs.FileVersionList, ipfsVersion)
			newIPFS = ipfs
		}
	}

	if (!prevVersionExists){
		newIPFS = IPFS{Id: uuid.New(), FileName: nonVersionedFileName, FileVersionList: []IPFSVersion{ipfsVersion}}
	}
	return newIPFS
}

func (ipfs *IPFS)GetIPFSJSON() string{
	ipfsJSON, err := json.Marshal(ipfs)
	if(err!=nil){
		fmt.Println("Unable to convert peer node ipfs list to json")
	}
	return string(ipfsJSON)
}


// This will create a new ipfs list
func NewIPFSList() IPFSList {
	ipfsMap := make(map[uuid.UUID]IPFS)
	ipfsList := IPFSList{IPFSMap: ipfsMap, UpdatedTime: time.Now()}
	return ipfsList
}


// This is called initially when node is started
// This will update all files present in IPFS directory to Peer IPFS List
func (ipfsList *IPFSList)FetchNodeIPFSList(IPFS_DIR string, peerNode Peer){
	ipfsList.mux.Lock()
	defer ipfsList.mux.Unlock()

	// This will create IPFS directory if it does not exist
	if _, err := os.Stat(IPFS_DIR); os.IsNotExist(err) {
		os.Mkdir(IPFS_DIR, 0700)
	}

	files, _ := ioutil.ReadDir(IPFS_DIR)
	for _, file := range files {
		fileName := file.Name()
		absoluteFilePath := path.Join(IPFS_DIR, fileName)
		fileStat, _ := os.Stat(absoluteFilePath)
		if(fileStat.Mode().IsRegular()){
			ipfs := NewIPFS(IPFS_DIR, file, peerNode, ipfsList)
			ipfsList.IPFSMap[ipfs.Id] = ipfs
		}
	}
}

// This will update the node ipfs list periodically
// Checks for newly created/modified files and then creates a new entry to IPFS list
func (ipfsList *IPFSList)PollNodeIPFSList(IPFS_DIR string, peerNode Peer) []IPFS{
	ipfsList.mux.Lock()
	defer ipfsList.mux.Unlock()
	RandomSleep()
	files, _ := ioutil.ReadDir(IPFS_DIR)
	newIPFSList := []IPFS{}
	for _, file := range files {
		fileName := file.Name()
		absoluteFilePath := path.Join(IPFS_DIR, fileName)
		fileStat, _ := os.Stat(absoluteFilePath)
		mtime := fileStat.ModTime()
		if(mtime.After(ipfsList.UpdatedTime) && fileStat.Mode().IsRegular()){
			ipfs := NewIPFS(IPFS_DIR, file, peerNode, ipfsList)
			fmt.Println("New IPFS File Found", ipfs.FileName, "at", ipfsList.UpdatedTime)
			newIPFSList = append(newIPFSList, ipfs)
			ipfsList.IPFSMap[ipfs.Id] = ipfs
		}
	}
	fmt.Println("All files have been scanned at", time.Now())
	ipfsList.UpdatedTime = time.Now()
	return newIPFSList
}

// This will sync ipfs heart beat list with node ipfs list
func (ipfsList *IPFSList)SyncNodeIPFSList(ipfsListJSON string){
	ipfsList.mux.Lock()
	defer ipfsList.mux.Unlock()
	var newIPFSList []IPFS
	err := json.Unmarshal([]byte(ipfsListJSON), &newIPFSList)
	if(err!=nil){
		fmt.Println(err)
	}
	for _, newIPFS := range newIPFSList {
		// TODO: Selectively merge ipfs list rather than overriding
		ipfsList.IPFSMap[newIPFS.Id] = newIPFS
	}
}


// This will get ipfs file available at node
// if ipfs file is not available at node download ipfs list from all peers and sync ipfs list
// then check if ipfs file is available at node
// if it available check if version is available to peer node
// if version is not available to peer node send forbidden request status code
func (ipfsList *IPFSList)GetIPFSFile(peerNode Peer, peerList PeerList, ipfsId uuid.UUID, versionId int) (IPFS, error) {
	var ipfsFile IPFS
	if ipfs, ok := ipfsList.IPFSMap[ipfsId]; ok {
		for _, version := range ipfs.FileVersionList{
			if(version.Id==versionId){
				for _, user :=  range version.VersionOwners {
					if (user.PeerNode.PeerId == peerNode.PeerId){
						return ipfs, nil
					}
				}
				return ipfsFile, errors.New("This peer is not authorized to access this version")
			}
		}
		return ipfsFile, errors.New("This version is not available in node ipfs list")
	}
	return ipfsFile, errors.New("This file is not available in node ipfs list")
}

// This will download ipfs from all peers and update node ipfs list
func (ipfsList *IPFSList)DownloadIPFSList(peerList PeerList) {
	for _, peerNode := range peerList.PeerMap {
		IPFS_LIST_URL := "http://" + peerNode.Address + "/ipfs/"
		fmt.Println("Initiating connection to Peer Server to update ipfs list at : ", IPFS_LIST_URL)
		res, err := http.Get(IPFS_LIST_URL)
		if (err != nil) {
			fmt.Println("Unable to fetch peers to find block at given height with hash")
		}else{
			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			blockBufferString := buf.String()
			newIPFSList := GetNodeIPFSListFromJSON(blockBufferString)
			ipfsList.mux.Lock()
			for ipfsId, ipfs := range(newIPFSList.IPFSMap){
				// TODO: Selectively merge ipfs list rather than overriding
				ipfsList.IPFSMap[ipfsId]=ipfs
			}
			ipfsList.mux.Unlock()
		}
	}
}



// This will return the ipfs list available at node as json string
func (ipfsList *IPFSList)GetNodeIPFSJSON() string{
	ipfsListJSON, err := json.Marshal(ipfsList)
	if(err!=nil){
		fmt.Println("Unable to convert peer node ipfs list to json")
	}
	return string(ipfsListJSON)
}


// This will return the ipfs list available at node as json string
func GetNodeIPFSListFromJSON(ipfsListJSON string) IPFSList{
	var ipfsList IPFSList
	err := json.Unmarshal([]byte(ipfsListJSON), &ipfsList)
	if (err != nil) {
		fmt.Println("Unable to decode ipfslist json to ipfslist", err)
	}
	return ipfsList
}

// This will return the IPFS List of peer as json
func (ipfsList *IPFSList)ShowNodeIPFSList() string{
	return ipfsList.GetNodeIPFSJSON()
}