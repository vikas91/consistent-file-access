package models

import "sync"

type Transaction struct {
	Hash string
	PreviousHash string
	Status string
	BlockId int32
	Sender Peer
	Receiver Peer
	Type string // 0: Share, 1: Seed
	Fee float32
	TimeStamp int64
	mux sync.Mutex
}
