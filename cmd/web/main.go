package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"sabiraliyev.net/snippetbox/pkg/models"
	"time"

	"sabiraliyev.net/snippetbox/pkg/models/postgresql"

	_ "github.com/lib/pq"
)

type contextKey string

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass"
	dbname   = "snippetbox"
)
const contextKeyAccount = contextKey("account")
const contextKeyIsAuthenticated = contextKey("isAuthenticated")
const contextKeyIsAdministrator = contextKey("isAdministrator")

// Define an application struct to hold the application wide dependencies for the web application.
// For now we`ll only include fields for the two custom loggers, but we`ll add more to it as build process.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	messages interface {
		Insert(int, string, string) (int, error)
		Get(int) (*models.Message, error)
		Latest(int) ([]*models.Message, error)
		Delete(int) (*models.Message, error)
		Update(int, string, bool) (int, error)
	}
	templateCache map[string]*template.Template
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in hte addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a command-line flag for the MySQL DSN string.
	// dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Define a command-line flag for the session secret (a random key which will be used to encrypt
	// and authenticate session cookies). It should be 32 bytes long.
	secret := flag.String("secret", "i334343g3+g3gk3@i3g335+35g4589gj", "Secret key")

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

	// To keep the Main() function tidy, we put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(psqlInfo)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Defer a call to db.Close(), so that the connection pool is closed before the main() function exits.
	defer db.Close()

	err = db.Ping()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the session.New() function to initialize a new session manager, passing in the secret key
	// as the parameter. Then we configure it so session always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true // Set the Secure flag on the session cookies.

	// Setting SameSite=Strict will block the session cookie being send by user`s browser for ALL cross-site usage.
	// This includes when a user clicks on an external link to the application, meaning that after clicking the link
	// they will initially be treated as 'not logged in' even if they an active session containing their
	// "authenticatedUserID" value:
	// session.SameSite = http.SameSiteStrictMode // Default value is SameSite=Lax

	// Initialize an instance of application struct containing the dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &postgresql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &postgresql.UserModel{DB: db},
		messages:      &postgresql.MessageModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default LTS settings we want server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialise a http.Server struct. We set the Addr and Handler fields so that
	// the server uses the same network address and rotes as before, and set the ErrorLog field
	// so that the server now uses the custom errorlog logger in the event of any problem.
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(), // Call the app.routes() method.
		TLSConfig:    tlsConfig,    // Set the server`s TLSConfig field.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Write messages using two loggers.
	infoLog.Printf("Starting server on %s", *addr)
	// Use yhe ListenAndServeTLS() method to start the HTTPS server. We pass in the paths
	// to the TLS certificates and corresponding private key as the two parameters.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pull for a given DSN.
func openDB(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
