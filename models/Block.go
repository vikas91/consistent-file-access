package models

import "sync"

type Header struct {
	Height int32
	TimeStamp int64
	Hash string
	ParentHash string
	Size int32
	Nonce string
}


type Block struct {
	BlockId int32
	Header Header
	Value MerklePatriciaTrie
	mux sync.Mutex
}


type BlockJSON struct{

}