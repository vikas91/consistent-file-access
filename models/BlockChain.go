package models

import "sync"

type BlockChain struct {
	Chain map[int32][]Block
	Length int32
	Difficulty int8
	mux sync.Mutex
}

