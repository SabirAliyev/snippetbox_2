package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

// Define an application struct to hold the application wide dependencies for the web application.
// For now we`ll only include fields for the two custom loggers, but we`ll add more to it as build process.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

// The ServerError helper writes an error message and stack trace to the errorLog,
// then sends e generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// To report the file name and line number *one step back* in the stack trace.
	// We do this by setting the frame depth to 2.
	app.errorLog.Println(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper send a specific status code and corresponding description
// to the user. We`ll use this later in the book to send responses like 400 "Bad Request"
// when there`s a problem with the request that the user send.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we`ll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which send a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

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

	// Initialize a new instance of application struct containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialise a http.Server struct. We set the Addr and Handler fields so that
	// the server uses the same network address and rotes as before, and set the ErrorLog field
	// so that the server now uses the custom errorlog logger in the event of any problem.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Call the new app.routes() method.
	}

	// Write messages using two loggers, instead of standard logger (see above).
	infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(":4000", mux) // DEPRECATED
	// Call the ListenAnServe() method on our new http.Server struct.
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
