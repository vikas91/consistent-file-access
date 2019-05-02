package models

const IPFS_FOLDER = "/tmp/ipfs"

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
	SeedCost int32
	SeedCount int32
	SeedEnabled bool
	VersionOwners []IPFSUser
	VersionSeeders []IPFSUser
	CreatedTimeStamp int64
}

type IPFSUser struct {
	PeerNode Peer
	PeerFileKey string
}


