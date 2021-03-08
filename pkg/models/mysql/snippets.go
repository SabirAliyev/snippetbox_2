package mysql

import (
	"database/sql"
	"errors"
	"log"
	// Import the models package we crated. You need to prefix this with
	// whatever module path you set up back in chapter 02.02 (Project Setup and
	// Enabling Modules) so that the import statement looks like this:
	// "{your-module-path}/pkg/models".
	"sabiraliyev.net/snippetbox/pkg/models"
)

// Define a SnippetModule type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// Write the SQL statement we want to execute. We split it over two lines
	// for readability (which is why it`s surrounded with backquotes instead
	// of normal double quotes).
	stmt := `
		INSERT INTO snippets 
			("title", content, created, expires) 
		VALUES 
			($1, $2, NOW(), NOW() + $3 * INTERVAL '1 DAY') 
		RETURNING snippetId;
			`

	result, err := m.DB.Prepare(stmt)
	if err != nil {
		log.Fatal(err)
	}

	// Use the Scan() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	var snippetId int
	err = result.QueryRow(title, content, expires).Scan(&snippetId)
	if err != nil {
		return 0, err
	}

	return snippetId, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `
		SELECT snippetId, title, content, created, expires 
		FROM snippets 
		WHERE expires > NOW() 
		AND snippetId = $1;
		`

	// Use the QueryRow() method on the connection pool to execute the SQL statement,
	// passing the untrusted id variable as the value for the placeholder parameter.
	// This returns a pointer to a sql.Row object which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding field
	// in the Snippet struct. Notice that the arguments to row.Scan are *pointers* to the place
	// you want to copy the data into, and the number of arguments must be exactly the same as
	// the number of columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return sl.ErrNoRows error. We use
		// the errors.IS() function check for that error  specifically, and return our own
		// models.ErrNoRecord error instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}

//This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `
		SELECT snippetId, title, content, created, expires 
		FROM snippets 
		WHERE expires > NOW() 
		ORDER BY created 
		DESC LIMIT 10;
		`

	// Use the Query() method on the connection pool to execute our SQL statement.
	// This returns a sql.Rows resultset containing the result of the query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is always property closed
	// before the Latest() method returns. This defer statement should come *after*
	// you check an error  from the Query() method. Otherwise, if Query() returns an error,
	// you`ll get a panic trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects.
	var snippets []*models.Snippet

	// Use rows.Next to iterate through the rows in the resultset. This prepares the first
	// (and then each subsequent) row to be acted on by the rows.Scan() method. If
	// iteration iteration over all the rows completes then the resultset automatically
	// closes itself and frees-up the underlying database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the new
		// Snippet object that we created. Again, the arguments to row.Scan() must be
		// pointers to the place you want to copy the data into, and the number of arguments
		// must be exactly the same as the number of columns returned by the statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// when the rows.Next() loop has finished we call.Err() to retrieve any error that was
	// encountered during the iteration. It`s important to call this - don`t assume that a
	// successful iteration was completed over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

// Mark snippet as Deleted. No actually removal will be performed.
func (m *SnippetModel) Delete() error {
	stmt := `
		UPDATE snippets 
		SET deleted = true 
		WHERE snippetId = $1;
		`

	_, err := m.DB.Query(stmt)

	return err
}
