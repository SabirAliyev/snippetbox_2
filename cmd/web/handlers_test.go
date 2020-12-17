package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a new instance of our application struct. For now, this just contains
	// a couple of mock loggers (which discard everything written to them.
	app := &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
	}

	// We than use the httptest.NewTLServer() function to create a new test server, passing
	// in the value returned by our app.routes() method as the handler for the server. This
	// starts up an HTTPS server which listens on a randomly-chosen port of your local machine
	// for the duration of the test. Notice that we defer a call to ts.Close() to shutdown the
	// server the test finishes.
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained in the ts.URL field.
	// We can use it along with the ts.client().Get()  method to make a GET /ping request against
	// the test server. This returns ah http.Response struct containing the response.
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	// We can then check the value of the response status code and body, using the same code as before.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// And we can check that the response body written by the ping handler equals OK.
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
