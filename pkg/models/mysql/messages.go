package mysql

import (
	"database/sql"
	"errors"
	"log"
	"sabiraliyev.net/snippetbox/pkg/models"
	"strconv"
)

type MessageModel struct {
	DB *sql.DB
}

func (m *MessageModel) Insert(userId, content, expires string) (int, error) {
	stmt := `INSERT INTO messages (userId, content, created, expires) VALUES($1, $2, NOW(), NOW() + 365 * INTERVAL '1 DAY') RETURNING messageId`

	result, err := m.DB.Prepare(stmt)
	if err != nil {
		log.Fatal(err)
	}

	var messageId int
	err = result.QueryRow(userId, content, expires).Scan(&messageId)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}

func (m *MessageModel) Get(id int) (*models.Message, error) {
	stmt := `SELECT messageId, userId, content, created, expires, deleted FROM snippets WHERE expires > NOW() AND id = $1`

	row := m.DB.QueryRow(stmt, id)
	msg := &models.Message{}
	err := row.Scan(&msg.MessageID, &msg.UserId, &msg.Content, &msg.Created, &msg.Expires, &msg.Deleted)
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
	stmt := `SELECT messageId, userId, content, created, expires FROM messages WHERE expires > NOW() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	message := []*models.Message{}

	for rows.Next() {
		msg := &models.Message{}
		err = rows.Scan(&msg.MessageID, &msg.UserId, &msg.Content, &msg.Created, &msg.Expires)
		if err != nil {
			return nil, err
		}
		message = append(message, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return message, nil
}

func (m *MessageModel) Delete(id int) (int, error) {
	idStr := strconv.Itoa(id)
	errorCode := 0
	stmt := "UPDATE messages SET deleted = true WHERE id =" + idStr + ";"

	_, err := m.DB.Prepare(stmt)
	if err != nil {
		log.Fatal(err)
	}

	_, err = m.DB.Exec(stmt)
	if err != nil {
		errorCode = 0
	}
	return errorCode, nil
}
