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

# Application Layer

 