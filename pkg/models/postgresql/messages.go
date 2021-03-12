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
		SELECT messageId, userId, userName, content, date, expires, edited, status 
		FROM messages 
		WHERE expires > NOW() 
		AND messageId = $1;
		`

	row := m.DB.QueryRow(stmt, id)
	msg := &models.Message{}
	err := row.Scan(&msg.MessageID, &msg.UserId, &msg.User, &msg.Content, &msg.Date, &msg.Expires, &msg.Edited, &msg.Status)
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
		SELECT messageId, users.name, content, date, expires, edited, status
		FROM messages 
		JOIN users ON messages.userid = users.userId 
		WHERE expires > NOW() 
		AND status != 2 
		LIMIT 100;
		`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{}

		err = rows.Scan(&msg.MessageID, &msg.User, &msg.Content, &msg.Date, &msg.Expires, &msg.Edited, &msg.Status)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MessageModel) Update(id int, content string) (int, error) {
	stmt := `
		UPDATE messages 
		SET content = $2 
		WHERE messageId =  $1;
		`

	result, err := m.DB.Prepare(stmt)
	if err != nil {
		log.Fatal(err)
	}
	var messageId int
	err = result.QueryRow(id, content).Scan(&messageId)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}

func (m *MessageModel) Delete(id int) (*models.Message, error) {
	stmt := `
		UPDATE messages 
		SET status = 2 
		WHERE messageId =  $1;
		`

	row := m.DB.QueryRow(stmt, id)
	msg := &models.Message{}
	err := row.Scan(&msg.MessageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return msg, nil
}
