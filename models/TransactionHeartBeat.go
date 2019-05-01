package models

type TransactionHeartBeat struct {
	TransactionId int32  `json:"transactionId"`
	TransactionJson string `json:"transactionJson"`
	PeerNode Peer `json:"peerNode"`
	Hops int32  `json:"hops"`
}
