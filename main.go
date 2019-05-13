package main

import (
	appRouter "github.com/vikas91/consistent-file-access/application/routers"
	appHandlers "github.com/vikas91/consistent-file-access/application/handlers"
	bcRouter "github.com/vikas91/consistent-file-access/block-chain/routers"
	bcHandlers "github.com/vikas91/consistent-file-access/block-chain/handlers"
	"log"
	"net/http"
	"os"
)

// This function will be executed before everything else.
// This will be used to read config parameters to start the node
func init() {

}

func main() {
	//TODO: Should read config parameters and run either application server or block node
	if len(os.Args) > 1 {
		bcHandlers.InitializePeerNode(os.Args)
		router := bcRouter.NewRouter()
		log.Fatal(http.ListenAndServe(":" + os.Args[1], router))
	} else {
		appHandlers.InitializeApplication()
		router := appRouter.NewRouter()
		log.Fatal(http.ListenAndServe(":6686", router))
	}
}
