package main

import (
	"github.com/golangcollege/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"sabiraliyev.net/snippetbox/pkg/models/mock"
	"testing"
	"time"
)

// The helper which returns an instance of our application struct, containing mock dependencies.
func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// Create a session manager instance with the same settings as production.
	session := sessions.New([]byte("ie489jn39kvr3jud0kn39wd7k3v5kw93"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Initialize the dependencies, using the mocks for the logger and database models.
	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		templateCache: templateCache,
		users:         &mock.UserModel{},
	}
}

// Define a custom testServer type which anonymously embeds an httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// The helper which  initialize and returns a new instance of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Initialize new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the cookie jar to the client, so that response cookies are stored and then sent
	// with subsequent requests.
	ts.Client().Jar = jar

	// Disable direct-following for the client. Essentially this function is called after
	// 3xx response is received by the client, and returning the http.ErrUserLastResponse
	// error forces it to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// The get method on our custom testServer type. This makes a GET request to a given url path
// on the test server, and returns the response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
