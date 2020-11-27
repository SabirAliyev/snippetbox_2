package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in hte addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This readr in the command-line flag value and assigns it to the addr variable.
	// You need to call it *before* you use the addr variable. Otherwise it will always
	// contain the default value of ":4000".
	// If any errors encountered during parsing the application will be terminated.
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for all URL
	// paths that start with "/static/". For matching paths, we strip  the "/static" prefix
	// before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e. prefix
	// it with the * symbol) before using it. Note that we`re using the log.Printf()
	// function to interpolate the address with the log message.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
