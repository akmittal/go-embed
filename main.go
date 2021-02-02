package main

import (
	"embed"
	"io/fs"
	"fmt"
	"net/http"
)

//go:embed client/build
var content embed.FS

func clientHandler() http.Handler {
	fsys := fs.FS(content)
	contentStatic, _ := fs.Sub(fsys, "client/build")
	return http.FileServer(http.FS(contentStatic))

}
func main() {
	port := ":3010"
	mux := http.NewServeMux()
	mux.Handle("/", clientHandler())
	fmt.Printf("Listening on port %v", port)
	http.ListenAndServe(port, mux)

}
