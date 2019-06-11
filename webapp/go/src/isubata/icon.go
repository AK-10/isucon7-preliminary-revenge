package main

import (
	"database/sql"
	"github.com/labstack/echo"
	"strings"
	"net/http"
	"github.com/gomodule/redigo/redis"
)

func getIcon(c echo.Context) error {
	var name string
	var data []byte

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	defer conn.Close()

	// まずキャッシュから取り出す
	
	name = c.Param("file_name")
	data, err = redis.Bytes(conn.Do("GET", name))
	if err != redis.ErrNil {
		err := db.QueryRow("SELECT name, data FROM image WHERE name = ?",
		c.Param("file_name")).Scan(&name, &data)
		if err == sql.ErrNoRows {
			return echo.ErrNotFound
		}
		if err != nil {
			return err
		}
		if err := setIcon(name, data); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	mime := ""
	switch true {
	case strings.HasSuffix(name, ".jpg"), strings.HasSuffix(name, ".jpeg"):
		mime = "image/jpeg"
	case strings.HasSuffix(name, ".png"):
		mime = "image/png"
	case strings.HasSuffix(name, ".gif"):
		mime = "image/gif"
	default:
		return echo.ErrNotFound
	}
	return c.Blob(http.StatusOK, mime, data)
}

func setIcon(name string, data []byte) error {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", name, string(data))
	if err != nil {
		return err
	}

	return nil
}