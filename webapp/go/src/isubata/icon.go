package main

import (
	"database/sql"
	"github.com/labstack/echo"
	"strings"
	"net/http"
)

var iconCache = map[string]([]byte){}

func getIcon(c echo.Context) error {
	var name string
	var data []byte
	// まずキャッシュから取り出す
	if len(iconCache[c.Param("file_name")]) > 0 {
		name = c.Param("file_name")
		data = iconCache[c.Param("file_name")]
	} else {
		err := db.QueryRow("SELECT name, data FROM image WHERE name = ?",
			c.Param("file_name")).Scan(&name, &data)
		if err == sql.ErrNoRows {
			return echo.ErrNotFound
		}
		if err != nil {
			return err
		}
		iconCache[name] = data
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