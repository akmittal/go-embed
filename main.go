package main

import (
	"embed"
	"io/fs"
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
	mux := http.NewServeMux()
	mux.Handle("/", clientHandler())
	http.ListenAndServe(":3000", mux)

}
