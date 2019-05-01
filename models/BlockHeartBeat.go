package models


type BlockHeartBeat struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	BlockId          int32  `json:"blockId"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}
