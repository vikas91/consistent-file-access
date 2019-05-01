package models

type IPFS struct {
	Id int32
	FileName string
	Address string
	FileVersion []IPFSVersion
}

type IPFSVersion struct {
	Id int32
	PreviousVersionHash string
	CurrentHash string
	VersionOwners []IPFSUser
	VersionSeeders []IPFSUser
	CreatedTimeStamp int64
}

type IPFSUser struct {
	PeerNode Peer
	PeerFileKey string
}