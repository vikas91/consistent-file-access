package models

type TransactionHeartBeat struct {
	TransactionId int32  `json:"transactionId"`
	TransactionJson string `json:"transactionJson"`
}

type SignedTransactionHeartBeat struct {
	PeerNodeList []Peer
	SignedTransactionHB string
	Hops int32  `json:"hops"`
}