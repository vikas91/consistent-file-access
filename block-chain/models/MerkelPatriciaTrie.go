package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
	"sync"
)

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


func NewMerklePatriciaTrie() MerklePatriciaTrie {
	mpt := MerklePatriciaTrie{db: make(map[string]Node), root: ""}
	return mpt
}

func (mpt *MerklePatriciaTrie) getNodeValue(node Node, hexArray []uint8, transactionMap map[string]string){
	if(node.node_type==1){
		for index, element := range node.branch_value {
			if(element!=""){
				// Adding the 16th branch value to transactiondb
				if(index==16){
					//fmt.Println("Transaction on branch at 16th index: ",  hexArray, compact_encode(hexArray), string(compact_encode(hexArray)[:]))
					transactionMap[string(compact_encode(hexArray)[1:])] = element
				}else{
					//fmt.Println("Transaction on branch, Current Hex: ",  hexArray)
					mpt.getNodeValue(mpt.db[element], append(hexArray, uint8(index)), transactionMap)
				}
			}
		}
	}else{
		hexArray = append(hexArray, compact_decode(node.flag_value.encoded_prefix)...)
		if(check_is_leaf(node.flag_value.encoded_prefix)){
			//fmt.Println("Transaction on leaf, Current Hex: ",  hexArray, compact_encode(hexArray), string(compact_encode(hexArray)), node.flag_value.value)
			transactionMap[string(compact_encode(hexArray)[1:])] = node.flag_value.value
		}else{
			//fmt.Println("Transaction on extension, Current Hex: ",  hexArray)
			mpt.getNodeValue(mpt.db[node.flag_value.value], hexArray, transactionMap)
		}
	}
}


// This function converts mpt into a map of key and values
func (mpt *MerklePatriciaTrie) GetTransactions() map[string]string{
	mpt_root := mpt.root
	transactionMap := make(map[string]string)
	if(mpt_root==""){
		return transactionMap
	}
	mpt.getNodeValue(mpt.db[mpt_root], []uint8{}, transactionMap)
	return transactionMap
}

func (mpt *MerklePatriciaTrie) find_node(hashString string, hex_array []uint8) (string, error) {
	mpt_node := mpt.db[hashString]
	//fmt.Println("Find called with", hex_array)
	//Recursive call for branch node
	if(mpt_node.node_type==1){
		//fmt.Println("Branch Node called")
		// Case for extension node points to branch node and hex_array is empty
		if(len(hex_array)==0){
			return mpt_node.branch_value[16], nil
		}
		// Case for branch node index is null but hex_array exists
		if(mpt_node.branch_value[hex_array[0]]=="" && len(hex_array)>0){
			return "", errors.New("path_not_found")
		}
		//fmt.Println("Index to find", hex_array[0])
		//fmt.Println("branch node", mpt_node.branch_value)
		return mpt.find_node(mpt_node.branch_value[hex_array[0]], hex_array[1:])
	}else{
		//fmt.Println("Find in Extension/Leaf node called")
		//fmt.Println(mpt_node)
		encoded_prefix := mpt_node.flag_value.encoded_prefix
		//fmt.Println(encoded_prefix, hex_array)
		// Case for leaf node with empty prefix
		if(len(compact_decode(encoded_prefix)) == 0 && len(hex_array)==0){
			return mpt_node.flag_value.value, nil
		}else{
			mpt_node_hex_array := compact_decode(encoded_prefix)
			if(check_is_leaf(encoded_prefix)){
				//fmt.Println("Find Leaf node called")
				// Case for leaf node without empty prefix
				if(reflect.DeepEqual(mpt_node_hex_array,hex_array)){
					return mpt_node.flag_value.value, nil
				}
				return "", errors.New("path_not_found")
			}else{
				//fmt.Println("Find in Extension node called")
				// Case for extension node with
				node_hex_length := len(mpt_node_hex_array)
				var i int;
				for i = 0; i < node_hex_length; i++ {
					// Case when hex array size is less than extension node size
					if(i>=len(hex_array)){
						return "", errors.New("path_not_found")
					}
					if(mpt_node_hex_array[i]!=hex_array[i]){
						break
					}
				}
				if(i<node_hex_length){
					return "", errors.New("path_not_found")
				}else{
					return mpt.find_node(mpt_node.flag_value.value, hex_array[i:])
				}
			}
		}
	}
	return "", errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	hex_array := compact_decode([]byte(key))
	if (mpt.root==""){
		return "", errors.New("path_not_found")
	}else{
		return mpt.find_node(mpt.root, hex_array)
	}
}

func (mpt *MerklePatriciaTrie) insert_node(mpt_node Node, hex_array []uint8, new_value string) string{
	current_node_hash := mpt_node.hash_node()
	//Insert onto branch node
	if(mpt_node.node_type==1){
		// Insert into branch node with empty hex array
		if(len(hex_array)==0){
			//fmt.Println("Insert into branch node with empty hex array", mpt_node, new_value)
			mpt_node.branch_value[16] = new_value
		}else{
			index := hex_array[0]
			// Insert into branch node without hashvalue at index
			if(mpt_node.branch_value[index]==""){
				//fmt.Println("Insert called on branch node without hashvalue at index", index, hex_array)
				hex_array = append(hex_array[1:], 16)
				flag_value := Flag_value{encoded_prefix: compact_encode(hex_array), value: new_value}
				new_node := Node{node_type: 2, flag_value: flag_value}
				new_hash := new_node.hash_node()
				mpt_node.branch_value[index] = new_hash
				mpt.db[new_hash] = new_node
			}else{
				// Insert into branch node with hashvalue at index
				//fmt.Println("Insert called on branch node with hashvalue at index")
				branch_node := mpt.db[mpt_node.branch_value[index]]
				new_hash := mpt.insert_node(branch_node, hex_array[1:], new_value)
				mpt_node.branch_value[index] = new_hash
			}
		}
	}else{
		// Insert into leaf or extension node
		encoded_prefix := mpt_node.flag_value.encoded_prefix
		mpt_node_hex_array := compact_decode(encoded_prefix)
		//fmt.Println("Insert called on leaf/extension node", mpt_node_hex_array, hex_array)
		// Case for hex array equal to node hex
		if(reflect.DeepEqual(mpt_node_hex_array, hex_array)){
			if(check_is_leaf(encoded_prefix)){
				//fmt.Println("Updating leaf node value")
				mpt_node.flag_value.value = new_value
			}else{
				//fmt.Println("Updating extension node value")
				branch_node := mpt.db[mpt_node.flag_value.value]
				branch_node_hash := mpt.insert_node(branch_node, []uint8{}, new_value)
				mpt_node.flag_value.value = branch_node_hash
			}
		}else{
			var i int;
			for i = 0; i < len(mpt_node_hex_array); i++ {
				if(i==len(hex_array) || mpt_node_hex_array[i]!=hex_array[i]){
					break
				}
			}
			prefix_first_split_array:= mpt_node_hex_array[:i]
			prefix_second_split_array := mpt_node_hex_array[i:]
			hex_second_split_array := hex_array[i:]
			//fmt.Println("Split hex pattern values", prefix_first_split_array, prefix_second_split_array, hex_second_split_array)

			// Convert Node to branch node
			if(len(prefix_first_split_array)==0) {
				//fmt.Println("Converting leaf/extension node into branch node", mpt_node)
				mpt_node.node_type = 1
				mpt_flag_value := mpt_node.flag_value.value
				mpt_node.flag_value.encoded_prefix = []byte("")
				mpt_node.flag_value.value = ""

				if(check_is_leaf(encoded_prefix)) {
					if(len(prefix_second_split_array) > 1) {
						//fmt.Println("Converting leaf node into branch node")
						prefix_flag_value := Flag_value{encoded_prefix: compact_encode(append(prefix_second_split_array[1:], 16)), value: mpt_flag_value}
						prefix_leaf_node := Node{node_type: 2, flag_value: prefix_flag_value}
						prefix_leaf_hash := prefix_leaf_node.hash_node()
						mpt_node.branch_value[prefix_second_split_array[0]] = prefix_leaf_hash
						mpt.db[prefix_leaf_hash] = prefix_leaf_node
					}else{
						mpt_node.branch_value[16] = mpt_flag_value
					}

				} else {
					if(len(prefix_second_split_array) > 1) {
						//fmt.Println("Converting extension node into branch node and adding extension node")
						prefix_flag_value := Flag_value{encoded_prefix: compact_encode(prefix_second_split_array[1:]), value: mpt_flag_value}
						prefix_leaf_node := Node{node_type: 2, flag_value: prefix_flag_value}
						prefix_leaf_hash := prefix_leaf_node.hash_node()
						mpt_node.branch_value[prefix_second_split_array[0]] = prefix_leaf_hash
						mpt.db[prefix_leaf_hash] = prefix_leaf_node
					} else {
						//fmt.Println("Converting extension node into branch node without adding extension node")
						mpt_node.branch_value[prefix_second_split_array[0]] = mpt_flag_value
					}
				}
				hex_flag_value := Flag_value{encoded_prefix: compact_encode(append(hex_second_split_array[1:], 16)), value: new_value}
				hex_leaf_node := Node{node_type: 2, flag_value: hex_flag_value}
				hex_leaf_hash := hex_leaf_node.hash_node()
				mpt_node.branch_value[hex_second_split_array[0]] = hex_leaf_hash
				mpt.db[hex_leaf_hash] = hex_leaf_node
				fmt.Println(mpt_node)
			}else if((len(prefix_first_split_array) == len(mpt_node_hex_array)) && !check_is_leaf(encoded_prefix)){
				//fmt.Println("Matched all contents of extension Node. Passing to next node")
				mpt_node.flag_value.value = mpt.insert_node(mpt.db[mpt_node.flag_value.value], hex_second_split_array, new_value)
			}else{
				//fmt.Println("Converting leaf/extension node to extension node")
				if(check_is_leaf(prefix_first_split_array)){
					prefix_first_split_array = prefix_first_split_array[:len(prefix_first_split_array)-1]
				}

				new_branch_node := Node{node_type: 1}

				//fmt.Println("hex insert", new_branch_node)
				new_branch_node_hash := mpt.insert_node(new_branch_node, hex_second_split_array, new_value)

				//fmt.Println("Prefix Insert", mpt.db[new_branch_node_hash])
				new_branch_node_hash = mpt.insert_node(mpt.db[new_branch_node_hash], prefix_second_split_array, mpt_node.flag_value.value)


				mpt_node.flag_value.encoded_prefix = compact_encode(prefix_first_split_array)
				new_branch_node = mpt.db[new_branch_node_hash]
				mpt_node.flag_value.value = new_branch_node.hash_node()
			}
		}
	}

	mpt_new_hash := mpt_node.hash_node()
	mpt.db[mpt_new_hash] = mpt_node
	delete(mpt.db, current_node_hash)
	return mpt_new_hash
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	hex_array := compact_decode([]byte(key))
	//fmt.Println("Inserting hex with value", hex_array, new_value)
	// Insert for case of root is empty
	if (mpt.root==""){
		//fmt.Println("Adding leaf node in empty tree")
		hex_array = append(hex_array, 16)
		flag_value := Flag_value{encoded_prefix: compact_encode(hex_array), value: new_value}
		root_node := Node{node_type: 2, flag_value: flag_value}
		root_hash := root_node.hash_node()
		mpt.db[root_hash] = root_node
		mpt.root = root_hash
	}else{
		root_hash := mpt.insert_node(mpt.db[mpt.root], hex_array, new_value)
		mpt.root = root_hash
	}
	//fmt.Println(mpt.db)
	//fmt.Println("Root:", mpt.root)
}


func (mpt *MerklePatriciaTrie) delete_node(mpt_node Node, hex_array []uint8) (string, error) {
	//fmt.Println("Delete called on", mpt_node, hex_array)
	current_hash := mpt_node.hash_node()

	// Delete on leaf or extension node
	if (mpt_node.node_type == 2) {
		encoded_prefix := mpt_node.flag_value.encoded_prefix
		mpt_node_hex_array := compact_decode(encoded_prefix)
		if (check_is_leaf(mpt_node.flag_value.encoded_prefix)) {
			//fmt.Println("Delete called on leaf node")
			if (reflect.DeepEqual(mpt_node_hex_array, hex_array)) {
				//fmt.Println("Found Leaf Node with given path.")
				delete(mpt.db, current_hash)
				return "", nil
			} else {
				//fmt.Println("Leaf Node with given path does not exist.")
				return "", errors.New("path_not_found")
			}
		} else {
			//fmt.Println("Delete called on extension node")
			var i int;
			for i = 0; i < len(mpt_node_hex_array); i++ {
				if (i == len(hex_array) || mpt_node_hex_array[i] != hex_array[i]) {
					break
				}
			}

			child_hash := mpt_node.flag_value.value
			child_node := mpt.db[child_hash]
			if (i==len(mpt_node_hex_array)) {
				var new_child_hash string
				var error error
				if(i<=len(hex_array)){
					//fmt.Println("Found Extension Node exceeding given path")
					new_child_hash, error = mpt.delete_node(child_node, hex_array[i:])
					//fmt.Println("Extension node recursively changed", new_child_hash, error)
				}else{
					//fmt.Println("Found Extension Node matching given path")
					new_child_hash, error = mpt.delete_node(child_node, []uint8{})
				}

				if(error==nil){
					child_node = mpt.db[new_child_hash]
					//fmt.Println(child_node)
					if(check_is_leaf(child_node.flag_value.encoded_prefix)){
						//fmt.Println("Merging extension node into leaf after deletion")
						mpt_node_hex_array = append(mpt_node_hex_array, compact_decode(child_node.flag_value.encoded_prefix)...)
						mpt_node_hex_array = append(mpt_node_hex_array, 16)
						mpt_node.flag_value.encoded_prefix = compact_encode(mpt_node_hex_array)
						mpt_node.flag_value.value = child_node.flag_value.value
						delete(mpt.db, new_child_hash)
					}
					mpt_node_new_hash := mpt_node.hash_node()
					mpt.db[mpt_node_new_hash] = mpt_node
					delete(mpt.db, current_hash)
					return mpt_node_new_hash, nil
				}
				return "", errors.New("path_not_found")
			}
			//fmt.Println("Extension Node with given path does not exist.")
			return "", errors.New("path_not_found")
		}

	}else {
		//fmt.Println("Delete called on branch node")
		//fmt.Println(mpt_node, hex_array)
		if (len(hex_array) == 0) {
			//fmt.Println("Deleting value of branch node")
			mpt_node.branch_value[16] = ""
		} else if (mpt_node.branch_value[hex_array[0]] == "") {
			//fmt.Println("Cannot delete at given path in branch node")
			return "", errors.New("path_not_found")
		} else {
			child_hex := mpt_node.branch_value[hex_array[0]]
			//fmt.Println("Trying to delete child hex at branch node", child_hex)
			new_child_hex, error := mpt.delete_node(mpt.db[child_hex], hex_array[1:])
			//fmt.Println(new_child_hex, error)
			if (error == nil) {
				mpt_node.branch_value[hex_array[0]] = new_child_hex
			} else {
				return "", errors.New("path_not_found")
			}
		}
		non_empty_indexes := []int{}
		for index, element := range (mpt_node.branch_value) {
			if (element != "") {
				non_empty_indexes = append(non_empty_indexes, index)
			}
		}
		//fmt.Println(len(non_empty_indexes))
		if (len(non_empty_indexes) == 0) {
			//fmt.Println("Deleting on branch node now contains no elements")
			delete(mpt.db, current_hash)
			return "", nil
		} else if (len(non_empty_indexes) == 1) {
			//fmt.Println("Merging branch node to its child as its contains only one element")
			child_node_hash := mpt_node.branch_value[non_empty_indexes[0]]
			child_node := mpt.db[child_node_hash]
			child_hex := compact_decode(child_node.flag_value.encoded_prefix)
			child_hex = append([]uint8{uint8(non_empty_indexes[0])}, child_hex...)
			//fmt.Println("Merged child hex", child_hex)
			if (check_is_leaf(child_node.flag_value.encoded_prefix)) {
				child_hex = append(child_hex, 16)
			}
			//fmt.Println("Merged child hex", child_hex)
			child_node.flag_value.encoded_prefix = compact_encode(child_hex)
			child_node_hash = child_node.hash_node()
			mpt.db[child_node_hash] = child_node
			//fmt.Println(child_node_hash, child_node)
			delete(mpt.db, current_hash)
			return child_node_hash, nil
		}
		//fmt.Println("Updating non empty branch Node after deleting children")
		mpt_node_new_hash := mpt_node.hash_node()
		mpt.db[mpt_node_new_hash] = mpt_node
		delete(mpt.db, current_hash)
		return mpt_node_new_hash, nil
	}
}


func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	hex_array := compact_decode([]byte(key))
	if (mpt.root==""){
		//fmt.Println("Deleting on empty root")
		return "", errors.New("path_not_found")
	}
	// Deletes node if it exists in the path and updates the hash in mpt db
	mpt_root_hash := mpt.root
	//fmt.Println("Root Hash Before Delete", mpt.root)
	mpt_root_node := mpt.db[mpt.root]
	mpt_new_root_hash, error := mpt.delete_node(mpt_root_node, hex_array)
	if(error==nil){
		mpt.root = mpt_new_root_hash
	}
	fmt.Println(mpt_root_hash, mpt_new_root_hash, error)
	return mpt_root_hash, error
}

// This function takes hex array as input and generates an ascii array
func compact_encode(hex_array []uint8) []uint8 {
	hex_length := len(hex_array)
	//if(hex_length==0){
	//	return []uint8{}
	//}
	var term uint8;
	if(hex_array[hex_length-1]==16){
		term = 1
		hex_array = hex_array[:hex_length-1]
	}else{
		term = 0
	}

	var oddlen uint8 = uint8(len(hex_array) % 2)
	var flags uint8 = 2 * term + oddlen

	if(oddlen==1){
		hex_array = append([]uint8{flags}, hex_array...)
	}else{
		hex_array = append([]uint8{flags, 0}, hex_array...)
	}
	var ascii_array []uint8;
	for i := 0; i < len(hex_array); i=i+2 {
		ascii_array = append(ascii_array, 16*hex_array[i] + hex_array[i+1])
	}
	return ascii_array
}

func check_is_leaf(encoded_arr []uint8) bool {
	var hexArray []uint8;
	var nodeType bool;
	if(len(encoded_arr)==0){
		return  false
	}
	for _, element := range encoded_arr {
		hexArray = append(hexArray, element/16, element%16)
	}
	//fmt.Println("leaf Check", hexArray)
	if(hexArray[0]==3 || hexArray[0] ==2){
		nodeType = true
	}
	return nodeType
}

// This function takes ascii array as input and generates a hex array
func compact_decode(encoded_arr []uint8) []uint8 {
	if(len(encoded_arr)==0){
		return []uint8{}
	}

	var hexArray []uint8;
	for _, element := range encoded_arr {
		hexArray = append(hexArray, element/16, element%16)
	}
	//Adjust Flags related bytes
	if(hexArray[0]==1 || hexArray[0]==3){
		hexArray = hexArray[1:]
	}else if((hexArray[0]==0 && hexArray[1]==0) || (hexArray[0]==2 && hexArray[1]==0)){
		hexArray = hexArray[2:]
	}
	return hexArray
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = string(node.flag_value.encoded_prefix) + node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
	mpt.root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0] / 16 < 2
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart" + cur_hash + "HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart" + cur_hash + "HashEnd", fmt.Sprintf("Hash%v", i),  -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}