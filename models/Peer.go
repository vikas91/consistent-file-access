package models


type Peer struct {
	PeerId int32
	Address string
	Balance float32
	PublicKey string
}

type PeerTransactionList struct {
	PeerNode Peer
	TransactionMap map[string]Transaction
}

type PeerIPFSList struct{
	PeerNode Peer
	IPFSMap map[int32]IPFS
}


