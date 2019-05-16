# API Description for Decentralised Consistent File Access

# Block Chain Layer
<details>
<summary>GET /</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Get Node Details<br/>
<pre>
{
    "PeerId": uuid, 
    "Address": string, 
    "Balance": float32,		
    "PublicKey": rsa.PublicKey, 
} 
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>POST / </summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Updates the key pair of a node<br/>
<pre>
{
    "PeerId": uuid, 
    "Address": string, 
    "Balance": float32,		
    "PublicKey": rsa.PublicKey, 
} 
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /peers/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Get List of Peers available to peer Details<br/>
<pre>
[
	{
		"PeerId": uuid, 
		"Address": string, 
		"Balance": float32,		
		"PublicKey": rsa.PublicKey, 
	} 
]
</pre></td><td>Implemented</td></tr>
</table>
</details>


<details>
<summary>GET /start/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Starts Polling new files in IPFS directory. Registers Node on Application Server. Return the peer node details<br/>
<pre>
{
    "PeerId": uuid, 
    "Address": string, 
    "Balance": float32,		
    "PublicKey": rsa.PublicKey, 
} 
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /stop/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Stops Polling new files in IPFS directory.<br/>
<pre>
{
    "PeerId": uuid, 
    "Address": string, 
    "Balance": float32,		
    "PublicKey": rsa.PublicKey, 
} 
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /blockchain/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Shows the complete block chain.<br/>
<pre>

</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /blockchain/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Updates the block chain by polling blocks from peers<br/>
<pre>

</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /blockChain/canonical</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Shows the canonical block chain list<br/>
<pre>

</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /blockChain/transactions</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Shows the block chain transaction list<br/>
<pre>

</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /blockChain/canonical/transactions</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Shows the block chain canonical transaction list<br/>
<pre>

</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /block/</summary>
Responses:
<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Requests block from peers<br/>
<pre>
</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /block/heartbeat/receive</summary>
Responses:
<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Requests block heart beat from peers<br/>
<pre>
</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /transactions/</summary>
Responses:
<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Get the list of transactions available at a peer<br/>
<pre>
</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /transactions/</summary>
Responses:
<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Updates the list of transactions available at a peer<br/>
<pre>
</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /transactions/heartbeat/receive</summary>
Responses:
<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Receives the transaction heart beat received from peer node.<br/>
<pre>
</pre></td><td>To Be Implemented</td></tr>
</table>
</details>

# Application Layer
<details>
<summary>GET /ipfs/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Shows List of IPFS files information at Node.<br/>
<pre>
{
    "IPFSMap": {
        IPFSId<uuid>: {
            Id: uuid,
            FileName : string,
            FileVersionList : [
                Id: int32,
                PreviousVersionHash: string,
                CurrentVersionHash: string,
                SeedCost: float32,
                SeedCount: int32,
                SeedEnabled: bool,
                VersionOwners: [
                    {
                        "PeerId": uuid, 
                        "Address": string, 
                        "Balance": float32,		
                        "PublicKey": rsa.PublicKey, 
                    }
                ]
                VersionSeeds: [
                    {
                        "PeerId": uuid, 
                        "Address": string, 
                        "Balance": float32,		
                        "PublicKey": rsa.PublicKey, 
                    }
                ]
            ]
            CreatedTime: time.Time
        }
    "UpdatedTime": time.Time  
} 
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /ipfs/download/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>Should Download ipfs files from peer and update node ipfs file list. Returns list of ipfs files available at node<br/>
<pre>
{
    "IPFSMap": {
        IPFSId<uuid>: {
            Id: uuid,
            FileName : string,
            FileVersionList : [
                Id: int32,
                PreviousVersionHash: string,
                CurrentVersionHash: string,
                SeedCost: float32,
                SeedCount: int32,
                SeedEnabled: bool,
                VersionOwners: [
                    {
                        "PeerId": uuid, 
                        "Address": string, 
                        "Balance": float32,		
                        "PublicKey": rsa.PublicKey, 
                    }
                ]
                VersionSeeds: [
                    {
                        "PeerId": uuid, 
                        "Address": string, 
                        "Balance": float32,		
                        "PublicKey": rsa.PublicKey, 
                    }
                ]
            ]
            CreatedTime: time.Time
        }
    "UpdatedTime": time.Time  
}  
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>POST ipfs/heartbeat/receive</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will be new heart beat received for IPFS list. Should check for signature and then add to node ipfs list. Forward Hear beat if hop count is greater than zero<br/>
<pre>
 
</pre></td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /ipfs/{ipfs_id}/version/{version_id}/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will retreive the IPFS file. Should retrieve file iff node is either in IPFS file owners, shared users or seeders<br/>
<pre>
{
    "IPFSData":"this is a dummy text to test file\n",
    "FileName":"test.txt"
}
</pre>
</td><td>Implemented upto owners</td></tr>
</table>
</details>

<details>
<summary>GET /ipfs/{ipfs_id}/versions</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will retreive the versions available to an IPFS entry<br/>
<pre>

</pre>
</td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /ipfs/{ipfs_id}/versions/{version_id}/owners/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will show the list of owners available to a ipfs file version <br/>
<pre>

</pre>
</td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /ipfs/{ipfs_id}/versions/{version_id}/seeds/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will show the list of seeds available to a ipfs file version <br/>
<pre>

</pre>
</td><td>To Be Implemented</td></tr>
</table>
</details>

<details>
<summary>POST /register/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will register the node on the application layer. <br/>
<pre>

</pre>
</td><td>Implemented</td></tr>
</table>
</details>

<details>
<summary>GET /nodes/</summary>

Responses:

<table>
	<tr><td>Code</td><td>Description</td><td>Status</td></tr>
	<tr><td>200</td><td>This will show the list of all node registered on the application layer. <br/>
<pre>

</pre>
</td><td>Implemented</td></tr>
</table>
</details>

 