package main

import (
	"errors"
	"fmt"
	"strconv"

	"net/http"
	"sabiraliyev.net/snippetbox/pkg/forms"
	"sabiraliyev.net/snippetbox/pkg/models"
)

// Change the signature of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual check of r.URL.PAth != "/" from this handler.

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the render helper.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Change the signature of the showSnippet() handler so it is defined as a method against *application.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn`t strip the colon from the names capture key,
	// so we need to get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	// Use the SnippetModel object`s Get method to retrieve the data for a specific record based on ID.
	// If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Pass the flash message to the template.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) showAdminPage(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "admin.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showChatPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyAccount).(*models.User)
	if user != nil {
		m, err := app.messages.Latest()
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.render(w, r, "chat.page.tmpl", &templateData{
			Messages: m,
		})
	} else {
		http.Redirect(w, r, fmt.Sprintf("/home"), http.StatusSeeOther)
	}
}

func (app *application) getUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyAccount).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// Add new createSnippetForm handler, which for now a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass the new empty form.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("content")
	form.MaxLength("content", 200)
	if !form.Valid() {
		app.render(w, r, "chat.page.tmpl", &templateData{Form: form})
	}
	userId := app.getUser(r).ID
	userName := app.getUser(r).Name

	_, err = app.messages.Insert(userId, userName, form.Get("content"))
	if err != nil {
		app.serverError(w, err)
		return
	} else {
		http.Redirect(w, r, fmt.Sprintf("/message/chat"), http.StatusSeeOther)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which add any data in POST request bodies to the r.PostForm map.
	// This also works in the same way for PUT and PATCH requests. If there any errors, we use
	// our app.Client helper to send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the form, we use the
	// validation methods to check the validation.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn`t valid, redisplay the template passing in the form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
	}

	// Because the form data (with type url.Values) has been anonymously embedded in the form.Form struct,
	// we can use the Get() method to retrieve the validated value for the particular form filed.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved successfully1") and the
	// corresponding key ("flash") to the session data. Note that if there`s no session for the current user
	// (or their session has expired) the new, empty session for them will automatically be created
	// by the session middleware.
	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) deleteSnippet(w http.ResponseWriter, r *http.Request) {
	//form := forms.New(r.Form.Get())
	fmt.Println("deleteSnippet method...")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate the form contents using the form helper.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("name", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	// Try to create new user record in the database. If the email already exist
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirm that their signup worker
	// and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they`re not, add a generic error message
	// to the form failures map and redisplay the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add the ID of the current user to the session, so that they are now 'logged in'.
	app.session.Put(r, "authenticatedUserID", id)

	// Add the Administrator bool value to the session.
	var isAdmin = app.isUserAdmin(id)
	app.session.Put(r, "isAdministrator", isAdmin)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) isUserAdmin(id int) bool {
	var isAdmin bool
	user, err := app.users.Get(id)

	if err != nil {
		app.errorLog.Fatal(err)
	} else {
		if user.Administrator {
			isAdmin = true
		}
	}
	return isAdmin
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// remove the authenticateUserID from the session data so that the user is 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that the`re benn logged out.
	app.session.Put(r, "flash", "You`ve been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
