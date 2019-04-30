# Decentralised Consistent File Access on BlockChain

## Problem Statement:
We want a system which can provide anonymity of transactions involving sensitive information sharing or agreements i.e record of transfer(ROT) between two or multiple parties. The record of transfer might be a document(s) proving the transfer of content or an agreement between parties involved in the transfer. Consensus should be achieved that the transfer has happened yet the contents of ROT should remain confidential to the world unless the parties involved wanted to reveal it. Parties involved in the transaction should be able to modify the contents of transaction but every change should be traceable and agreed upon by parties involved in the ROT. The parties involved in the transaction should be able to prove that the transaction content is authentic (i.e transfer content is immutable at a point in time). The record of transfer is stored on a decentralised file system and is accessible to only the parties involved in the transaction.

## Why Blockchain:
Blockchain chained transactions,  block creation algorithms and distributed consensus makes the contents of transaction immutable once added to the blockchain. This immutability can be used by users to record agreements on blockchain. We would like transactions to be recorded, but have control over the shared access and consistency during ROT versioning. We will use signed signatures to ensure authenticity of file creations, Gossip Protocol to transmit changes across the network and consensus models to ensure consistent shared access of ROT.     

## Problems Solved:
1. Parties involved in transfer should be anonymous to the world
2. Transfer/Agreements between parties should be recorded on blockchain which makes it immutable
3. Files describing the transfer should have versions and restrictive access to only parties involved in transfer
4. When transfer happens from A to B and then from B to C, A should see only version 1 while B and C should see version 2
5. File storage should be a secure decentralised storage  with authentication and authorisation  
6. File transfer content is immutable at a point in time.


## Architecture:

### IPFS
UserA who is the owner of FileA wants to broadcast that he is owner of ROT. UserA creates an IFPS entry and broadcasts his IPFS list to his peers. Eventually every one in the network will know that FileA can be accessed from UserA. Before adding the entry to UserA’s IFPS list we first create an AES encryption key for ROT which creates an encrypted text of the file content. The AES encryption key is encrypted further with the public key of UserA. We then apply a hash of the encrypted text of ROT i.e SHA256. The entry in the IFPS list hash mapped to node address of UserA. We will have UserA signed signature showing that he is indeed the owner of ROT.

### Accessing IPFS Entry
When UserB eventually receives an update on the ROT created by UserA i.e IPFS list of ROT gets updated he will be able to access the file using the hash provided in the IPFS entry. UserB sends his public key to UserA node. If UserB doesn’t have access to file(read/write) the file access API returns 403 forbidden. If hash doesn’t exists at UserA side then it returns a 404. IPFS file list will have metadata keys for every file entry. UserA should use his private key and decrypt  his metadata key and then encrypt it again with UserB’s public key and check if that exists in the list of metadata keys of the file. If it exists then it means UserB has access to ROT and the request returns a 200 along with the encrypted text of file. UserB on receiving the encrypted text will decrypt his metadata key of file and then decrypt the encrypted file to access the contents of file.

### Shared Access
Let us say UserA is the owner of file and want to share the agreement (ROT) with UserB. UserA creates the hash of the file described as part of transaction adding metadata key of userB to the file and adds it to the IFPS ledger and gossips the update to UserB with his signed signature. UserB now only has read access to ROT. When file shared access message is received on UserB side he sends an affirmation to UserA with his signed  signature agreeing on the content of ROT.  Note here that UserB has access to file version at the point the shared access transaction is created on UserA side. If UserA creates a version 2 of ROT before UserB sends affirmation with his signed signature then UserA cannot create a new transaction until UserB sends his affirmation on version 2. 

### Shared Access Transactions
Let us say User A(sender) and User B(receiver) agree on the ROT we now create a transaction on blockchain. Note here that at the end of agreement UserA and UserB are both owners of ROT i.e they have rights to access the content of ROT and can modify its contents. Now at the end of agreement UserA and UserB both have hash of files in their IPFS file list. UserA and UserB both can publish the transaction. Miners verify the transaction by checking the hash of file and public keys list on file from IPFS list and add the transaction to blockchain/
 
### Transactions Broadcast
Transaction that is now broadcasted to the network contains IPFS file hash , list of public keys of users that the file has access to and incentive value to verify the transaction. Miners verify the transaction which then adds the transaction to the block. Transaction verification will check if the hash of file broadcasted is present in IPFS file list of UserA and UserB and that both the parties are in fact owners of ROT.

### File Seeding Broadcast
Let say UserA/UserB who have made a ROT want to keep a back up copy of ROT. Either parties will broadcast a seed request which contains the hash of ROT, incentive to seed, time period for seeding.  UserC who wants to seed the ROT between UserA & UserB will accept the file transfer which will then receive the AES encrypted file along with the metadata keys of file.  User C will now have ROT in his IPFS file list but only has it as a seeder. UserA IPFS file list is now updated with the seeds list for the ROT.


### File Versioning Broadcast
Let us say UserA and UserB have access to ROT. UserB updates the file and want to broadcast the update. In order for UserB to add version 2 to IPFS list UserA has to agree that UserB has made a change to ROT. UserB creates an update file broadcast similar to create file broadcast. The update file broadcast will have new hash of file, new metadata keys of users having access 
to version 1 of file and UserB signed signature like create file broadcast message. User A will accept the new share access and broadcast the update message. When UserB receives the update from A we add the new version hash to IPFS list which has previous version hash pointing to version 1.
