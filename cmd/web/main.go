package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	// Use log.New() to create a logger with writing information messages. This takes
	// three parameters: the destination to write the log to (os.Stdout) , a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// filename and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	// (Use log.Llongfile to include full file path on log output).

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
	// log.Printf("Starting server on %s", *addr) // DEPRECATED

	// Write messages using two loggers, instead of standard logger (see above).
	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(":4000", mux)
	errorLog.Fatal(err)
}
