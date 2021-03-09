package postgresql

import (
	"database/sql"
	"errors"
	"log"
	"sabiraliyev.net/snippetbox/pkg/models"
)

type MessageModel struct {
	DB *sql.DB
}

func (m *MessageModel) Insert(userId int, User, content string) (int, error) {
	stmt := `
		INSERT INTO messages 
		    (userid, userName, content, date, expires) 
		VALUES 
			($1, $2, $3, NOW(), NOW() + 365 * INTERVAL '1 DAY') 
		RETURNING messageId;
		`

	result, err := m.DB.Prepare(stmt)
	if err != nil {
		log.Fatal(err)
	}
	var messageId int
	err = result.QueryRow(userId, User, content).Scan(&messageId)
	if err != nil {
		return 0, err
	}
	return messageId, nil
}

func (m *MessageModel) Get(id int) (*models.Message, error) {
	stmt := `
		SELECT messageId, userId, userName, content, date, expires, deleted 
		FROM messages 
		WHERE expires > NOW() 
		AND messageId = $1;
		`

	row := m.DB.QueryRow(stmt, id)
	msg := &models.Message{}
	err := row.Scan(&msg.MessageID, &msg.UserId, &msg.Content, &msg.Date, &msg.Expires, &msg.Deleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return msg, nil
}

func (m *MessageModel) Latest() ([]*models.Message, error) {
	stmt := `
		SELECT users.name, content, date, expires
		FROM messages 
		JOIN users ON messages.userid = users.userId 
		WHERE expires > NOW() 
		AND deleted = false LIMIT 100;
		`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}

		err = rows.Scan(&message.User, &message.Content, &message.Date, &message.Expires)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MessageModel) Delete() error {
	stmt := `
		UPDATE messages 
		SET deleted = true 
		WHERE messageId =  $1;
		`

	_, err := m.DB.Query(stmt)

	return err
}