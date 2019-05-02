package models

import "sync"

type BlockChain struct {
	Chain map[int32][]Block
	Length int32
	PeerNodeCount int32
	Difficulty int8
	mux sync.Mutex
}

// This function will return an instance of new block chain after instantiation
func NewBlockChain() BlockChain {
	blockChain := BlockChain{Chain: make(map[int32][]Block), Length: int32(0)}
	return blockChain
}
