package models

import (
	"crypto/rsa"
	"github.com/google/uuid"
)

type Node struct {
	PeerId uuid.UUID `peerId`
	Address string `address`
	Balance float32 `balance`
	PublicKey rsa.PublicKey `publicKey`
}

type UserList map[Node]uuid.UUID


