package mysql

import (
	"reflect"
	"sabiraliyev.net/snippetbox/pkg/models"
	"testing"
	"time"
)

func TestUserModelGet(t *testing.T) {
	// Skip the test if the `-short` flag is provided when running the test.
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name      string
		ID        int
		wantUser  *models.User
		wantError error
	}{
		{
			name: "Valid ID",
			ID:   1,
			wantUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			ID:        0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			ID:        2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// initialize a connection pool to our test database, and refer a call to the teardown function,
			// so it is always run immediately before this sub-test returns.
			db, teardown := newTestDB(t)
			defer teardown()

			// create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Get() method and check that the return value and error match
			// the expected values for the sub-test.
			user, err := m.Get(tt.ID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v; got %v", tt.wantUser, user)
			}
		})
	}
}
