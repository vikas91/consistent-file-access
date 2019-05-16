package models

import (
	"sync"
	"time"
)

type Transaction struct {
	Hash string
	FileVersionHash string
	FileVersion int32
	Status string
	BlockId int32
	Type string // 0: Create, 1: Share, 2: Seed
	Fee float32
	TimeStamp time.Time
	mux sync.Mutex
}


// This will always create new transaction for latest version of IPFS
// This will sign the transaction with the owner of latest version of IPFS
// This will broadcast the newly created transaction to all node peers
//func CreateTransaction(ipfs IPFS) Transaction {
//
//}

