package handlers

import (
	"github.com/google/uuid"
	"github.com/vikas91/consistent-file-access/application/models"
	"net/http"
	"path"
)


const TEMPLATE_DIR = "/Users/victor/workspace/go/src/github.com/vikas91/consistent-file-access/application/templates/"

var nodeList models.NodeList

func InitializeApplication(){
	nodeList = make(map[models.Node]uuid.UUID)
}

func Index(w http.ResponseWriter, r *http.Request) {
	p := path.Dir(TEMPLATE_DIR+"index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

func RegisterNode(w http.ResponseWriter, r *http.Request) {

}

func ShowNodeList(w http.ResponseWriter, r *http.Request) {
	peerJSON := nodeList.GetNodeListJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peerJSON))
}
