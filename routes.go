package consistent_file_access

import "net/http"
import "github.com/vikas91/consistent-file-access/handlers"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
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
		"RestartNode",
		"GET",
		"/restart",
		handlers.RestartNode,
	},
	Route{
		"ShowBlockChain",
		"GET",
		"blockChain/",
		handlers.ShowBlockChain,
	},
	Route{
		"UpdateBlockChain",
		"POST",
		"blockChain/update",
		handlers.UpdateBlockChain,
	},
	Route{
		"ShowCanonicalBlockChain",
		"GET",
		"blockChain/canonical",
		handlers.ShowCanonicalBlockChain,
	},
	Route{
		"ShowBlockChainTransactions",
		"GET",
		"blockChain/transactions",
		handlers.ShowBlockChainTransactions,
	},
	Route{
		"ShowCanonicalBlockChainTransactions",
		"GET",
		"blockChain/canonical/transactions",
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
		"block/heartbeat/receive",
		handlers.BlockHeartBeatReceive,
	},
	Route{
		"ShowTransactions",
		"GET",
		"transactions/",
		handlers.ShowTransactions,
	},
	Route{
		"UpdateTransactions",
		"POST",
		"transactions/update",
		handlers.UpdateTransactions,
	},
	Route{
		"TransactionHeartBeatReceive",
		"POST",
		"transactions/heartbeat/receive",
		handlers.TransactionHeartBeatReceive,
	},
}