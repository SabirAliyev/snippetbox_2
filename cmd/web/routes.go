package main

import (
	"github.com/bmizerany/pat"
	"net/http"

	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a http.Handler instead of *http.ServerMux.
func (app *application) routes() http.Handler {
	// The middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// The middleware chain containing the middleware specific to our dynamic application routes.
	// Using the noSurf middleware on all 'dynamic' routes with authenticate() and authenticateAsAdmin() middleware.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate, app.authenticateCurrentUser, app.authenticateAsAdmin)

	mux := pat.New()

	// Snippet
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/message/chat", dynamicMiddleware.ThenFunc(app.showChatPage))
	mux.Post("/message/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createMessage))
	mux.Get("/message/:id/delete", dynamicMiddleware.ThenFunc(app.deleteMessage))
	mux.Get("/snippet/admin", dynamicMiddleware.ThenFunc(app.showAdminPage))
	mux.Post("/snippet/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	// /Snippet

	// User session
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	// /User session

	// Test
	mux.Get("/ping", http.HandlerFunc(ping))
	// /Test

	mux.Get("/ping", http.HandlerFunc(ping))

	// Create a file server which serves files out of the "./ui/static" directory. Note that the path given
	// to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for all URL paths that start with
	// "/static/". For matching paths, we strip  the "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Return the 'standard' middleware chain followed by servemux.
	return standardMiddleware.Then(mux)
}
