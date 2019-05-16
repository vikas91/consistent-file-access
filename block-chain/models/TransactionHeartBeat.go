package models

type SignedTransactionHeartBeat struct {
	SignedNode Peer `json:"peer"`
	ForwardNode Peer `json:"peer"`
	SignedTransaction string `json:"signedTransaction"`
	Transaction string `json:"transaction"`
	Hops int32  `json:"hops"`
}