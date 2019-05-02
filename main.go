package main

import (
	"github.com/vikas91/consistent-file-access/routers"
	"log"
	"net/http"
	"os"
)

func main() {
	router := routers.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":" + os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}
}
