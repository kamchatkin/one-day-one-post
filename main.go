package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"one-day-one-post/db"
	modelPost "one-day-one-post/models/post"
	"one-day-one-post/static"
	"one-day-one-post/utils"
	"time"
)

func main() {
	defer db.Close()

	http.HandleFunc("/", index)
	http.HandleFunc("/post", post)
	http.HandleFunc("/style.css", style)
	http.HandleFunc("/style.css.map", styleMap)
	http.HandleFunc("/favicon.ico", faviconIco)
	http.HandleFunc("/favicon.png", favicon)

	http.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("404. Not Found"))
		w.WriteHeader(http.StatusNotFound)
	})

	_ = http.ListenAndServe(":8080", nil)
}

// index список постов
func index(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.DB.Query(modelPost.SqlSelect())
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var (
			id        int64
			createdAt string
			text      string
		)
		err = rows.Scan(&id, &createdAt, &text)
		if err != nil {
			panic(err)
		}

		textRunes := []rune(text)
		fmt.Printf("id: %d, created_at: %s, text: %s \n\n", id, createdAt, string(textRunes[0:utils.MaxRune(textRunes, 50)]))
	}

	_, _ = w.Write(static.IndexPage)
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if r.Method != "POST" {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":    false,
			"error": "Принимаем только POST",
		})
		return
	}

	post := modelPost.Post{}

	_ = json.NewDecoder(r.Body).Decode(&post)

	_, err := db.DB.Exec(post.SqlCreate())
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":    false,
			"error": fmt.Sprintf("Ошибка при добавлении поста в БД(%s)", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   true,
		"post": post.Text,
		"time": post.Time.Format(time.RFC3339),
	})
}

func style(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	_, _ = w.Write([]byte(static.StylePage))
}

func styleMap(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte(static.StyleMapPage))
}

func favicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "image/png")
	_, _ = w.Write(static.FaviconImage)
}

func faviconIco(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Location", "/favicon.png")
	w.WriteHeader(http.StatusPermanentRedirect)
	_, _ = w.Write([]byte{})
}
