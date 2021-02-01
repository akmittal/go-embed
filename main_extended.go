package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
)

//go:embed client/build
var f embed.FS

type User struct {
	Name  string `json:"name"`
	Email string `json:"email`
}

var (
	users []User
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/user", userHandler)
	mux.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
func userHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" { // Handle User Get
		json.NewEncoder(rw).Encode(users)
		return
	} else if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, "Error reading data", 400)
		}

		var u User
		err = json.Unmarshal(body, &u)
		if err != nil {
			http.Error(rw, "Not valid user data", 400)
		}
		users = append(users, u)
		return
	}

}
func rootHandler(rw http.ResponseWriter, req *http.Request) {
	upath := req.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		req.URL.Path = upath
	}
	upath = path.Clean(upath)
	fsys := fs.FS(f)
	fmt.Println(upath)
	contentStatic, _ := fs.Sub(fsys, "client/build")
	if _, err := contentStatic.Open(strings.TrimLeft(upath, "/")); err != nil { // If file not found server index/html from root
		req.URL.Path = "/"
	}
	http.FileServer(http.FS(contentStatic)).ServeHTTP(rw, req)
}
