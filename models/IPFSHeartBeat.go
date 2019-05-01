package models

type IPFSHeartBeat struct {
	IPFSId          int32  `json:"blockId"`
	IPFSJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}