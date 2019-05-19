package main

import (
	"time"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type Message struct {
	ID        int64     `db:"id"`
	ChannelID int64     `db:"channel_id"`
	UserID    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func addMessage(channelID, userID int64, content string) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO message (channel_id, user_id, content, created_at) VALUES (?, ?, ?, NOW())",
		channelID, userID, content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}


func queryMessages(chanID, lastID int64) ([]Message, error) {
	msgs := []Message{}
	err := db.Select(&msgs, "SELECT * FROM message WHERE id > ? AND channel_id = ? ORDER BY id DESC LIMIT 100",
		lastID, chanID)
	return msgs, err
}

func getMessage(c echo.Context) error {
	userID := sessUserID(c)
	if userID == 0 {
		return c.NoContent(http.StatusForbidden)
	}

	chanID, err := strconv.ParseInt(c.QueryParam("channel_id"), 10, 64)
	if err != nil {
		return err
	}
	lastID, err := strconv.ParseInt(c.QueryParam("last_message_id"), 10, 64)
	if err != nil {
		return err
	}

	messages, err := queryMessages(chanID, lastID)
	if err != nil {
		return err
	}

	response := make([]map[string]interface{}, 100)
	rows, err := db.Query("SELECT m.id, m.content, m.created_at, u.name, u.display_name, u.avatar_icon FROM message m INNER JOIN user u ON m.user_id = u.id WHERE m.id > ? AND m.channel_id = ? ORDER BY m.id DESC LIMIT 100", lastID, chanID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var message Message
		var user User
		err := rows.Scan(&message.ID, &message.Content, &message.CreatedAt, &user.Name, &user.DisplayName, &user.AvatarIcon)
		if err != nil {
			return err
		}
		r := makeJSONMessage(message, user)
		response = append([]map[string]interface{}{r}, response...)
	}

	if len(messages) > 0 {
		_, err := db.Exec("INSERT INTO haveread (user_id, channel_id, message_id, updated_at, created_at)"+
			" VALUES (?, ?, ?, NOW(), NOW())"+
			" ON DUPLICATE KEY UPDATE message_id = ?, updated_at = NOW()",
			userID, chanID, messages[0].ID, messages[0].ID)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, response)
}

func makeJSONMessage(m Message, u User) map[string]interface{} {
	r := make(map[string]interface{})
	r["id"] = m.ID
	r["user"] = u
	r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
	r["content"] = m.Content
	return r
}

// func getMessage(c echo.Context) error {
// 	userID := sessUserID(c)
// 	if userID == 0 {
// 		return c.NoContent(http.StatusForbidden)
// 	}

// 	chanID, err := strconv.ParseInt(c.QueryParam("channel_id"), 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	lastID, err := strconv.ParseInt(c.QueryParam("last_message_id"), 10, 64)
// 	if err != nil {
// 		return err
// 	}

// 	messages, err := queryMessages(chanID, lastID)
// 	if err != nil {
// 		return err
// 	}

// 	response := make([]map[string]interface{}, 0)
// 	for i := len(messages) - 1; i >= 0; i-- {
// 		m := messages[i]

// 		r, err := jsonifyMessage(m)
// 		if err != nil {
// 			return err
// 		}
// 		response = append(response, r)
// 	}

// 	if len(messages) > 0 {
// 		_, err := db.Exec("INSERT INTO haveread (user_id, channel_id, message_id, updated_at, created_at)"+
// 			" VALUES (?, ?, ?, NOW(), NOW())"+
// 			" ON DUPLICATE KEY UPDATE message_id = ?, updated_at = NOW()",
// 			userID, chanID, messages[0].ID, messages[0].ID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return c.JSON(http.StatusOK, response)
// }

func jsonifyMessage(m Message) (map[string]interface{}, error) {
	u := User{}
	err := db.Get(&u, "SELECT name, display_name, avatar_icon FROM user WHERE id = ?",
		m.UserID)
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})
	r["id"] = m.ID
	r["user"] = u
	r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
	r["content"] = m.Content
	return r, nil
}

// {
// 	"id": message.ID,
// 	"user": User{},
// 	"date": message.CreatedAt,
// 	"content": message.Content
// }

// "INSERT INTO haveread (user_id, channel_id, message_id, updated_at, created_at) VALUES (?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE message_id = ?, updated_at = NOW()"