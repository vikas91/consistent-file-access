package models

import "sync"

type Transaction struct {
	Hash string
	PreviousHash string
	Status string
	BlockId int32
	Sender string
	Receiver string
	Type string
	Fee float32
	TimeStamp int64
	mux sync.Mutex
}
