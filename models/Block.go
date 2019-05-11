package models

import "sync"

type Header struct {
	Height int32
	Hash string
	ParentHash string
	Size int32
	Nonce string
	RewardValue int32
	TimeStamp int64
}

type Block struct {
	BlockId int32
	Header Header
	Value MerklePatriciaTrie
	mux sync.Mutex
}


type BlockJSON struct{

}