package models

import "sync"

type Flag_value struct {
	encoded_prefix []uint8
	value string
}

type Node struct {
	node_type int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value Flag_value
}

type MerklePatriciaTrie struct {
	db map[string]Node
	root string
	mux sync.Mutex
}
