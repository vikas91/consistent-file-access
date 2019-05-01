package models


type TransactionHeartBeat struct {
	IfNewTransaction  bool   `json:"ifNewTransaction"`
	TransactionId          int32  `json:"transactionId"`
	TransactionJson   string `json:"transactionJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}
