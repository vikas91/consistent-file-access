package models

import "sync"

type Transaction struct {
	Hash string
	PreviousHash string
	Status string
	BlockId int32
	Sender string
	Receiver string
	Value float32
	Fee float32
	mux sync.Mutex
}
