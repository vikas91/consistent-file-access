package handlers

import (
	"net/http"
)




func Index(w http.ResponseWriter, r *http.Request) {
	//lp := filepath.Join("templates", "index.html")
	//fp := filepath.Join("templates", filepath.Clean(r.URL.Path))
	//
	//tmpl, _ := template.ParseFiles(lp, fp)
	//tmpl.ExecuteTemplate(w, "layout", nil)
}

func RegisterNode(w http.ResponseWriter, r *http.Request) {

}
