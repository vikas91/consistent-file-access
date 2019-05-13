package routers

import "net/http"
import "github.com/vikas91/consistent-file-access/block-chain/handlers"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"GetNodeDetails",
		"GET",
		"/",
		handlers.GetNodeDetails,
	},
	Route{
		"GetNodePeerList",
		"GET",
		"/peers/",
		handlers.GetNodePeerList,
	},
	Route{
		"UpdateNodeDetails",
		"POST",
		"/",
		handlers.UpdateNodeDetails,
	},
	Route{
		"StartNode",
		"GET",
		"/start",
		handlers.StartNode,
	},
	Route{
		"StopNode",
		"GET",
		"/stop",
		handlers.StopNode,
	},

	Route{
		"ShowBlockChain",
		"GET",
		"/blockChain/",
		handlers.ShowBlockChain,
	},
	Route{
		"UpdateBlockChain",
		"POST",
		"/blockChain/update",
		handlers.UpdateBlockChain,
	},
	Route{
		"ShowCanonicalBlockChain",
		"GET",
		"/blockChain/canonical",
		handlers.ShowCanonicalBlockChain,
	},
	Route{
		"ShowBlockChainTransactions",
		"GET",
		"/blockChain/transactions",
		handlers.ShowBlockChainTransactions,
	},
	Route{
		"ShowCanonicalBlockChainTransactions",
		"GET",
		"/blockChain/canonical/transactions",
		handlers.ShowCanonicalBlockChainTransactions,
	},
	Route{
		"RequestBlock",
		"POST",
		"/block/",
		handlers.RequestBlock,
	},
	Route{
		"BlockHeartBeatReceive",
		"POST",
		"/block/heartbeat/receive",
		handlers.BlockHeartBeatReceive,
	},
	Route{
		"ShowTransactions",
		"GET",
		"/transactions/",
		handlers.ShowTransactions,
	},
	Route{
		"UpdateTransactions",
		"POST",
		"/transactions/update",
		handlers.UpdateTransactions,
	},
	Route{
		"TransactionHeartBeatReceive",
		"POST",
		"/transactions/heartbeat/receive",
		handlers.TransactionHeartBeatReceive,
	},
	Route{
		"ShowIPFSList",
		"GET",
		"/ipfs/",
		handlers.ShowIPFSList,
	},
	Route{
		"UpdateIPFSList",
		"POST",
		"/ipfs/download",
		handlers.DownloadIPFSList,
	},
	Route{
		"IPFSHeartBeatReceive",
		"POST",
		"ipfs/heartbeat/receive",
		handlers.IPFSHeartBeatReceive,
	},
	Route{
		"ShowIPFSFile",
		"GET",
		"/ipfs/{ipfs_id}",
		handlers.ShowIPFSFile,
	},
	Route{
		"ShowIPFSFileVersions",
		"GET",
		"/ipfs/{ipfs_id}/versions",
		handlers.ShowIPFSFileVersions,
	},
	Route{
		"ShowIPFSSeedRequests",
		"GET",
		"/ipfs/seeds/",
		handlers.ShowIPFSSeedRequests,
	},
	Route{
		"ShareRequestIPFSFile",
		"POST",
		"/ipfs/{ipfs_id}/share",
		handlers.ShareRequestIPFSFile,
	},
	Route{
		"SeedRequestIPFSFile",
		"POST",
		"/ipfs/{ipfs_id}/seed",
		handlers.SeedRequestIPFSFile,
	},
}