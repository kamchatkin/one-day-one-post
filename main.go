package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"one-day-one-post/db"
	modelPost "one-day-one-post/models/post"
	"one-day-one-post/static"
	"strconv"
	"time"
)

func main() {
	defer db.Close()

	http.HandleFunc("/", index)
	http.HandleFunc("/api/posts", apiPosts)

	http.HandleFunc("/post", post)

	http.HandleFunc("/style.css", style)
	http.HandleFunc("/favicon.png", favicon)

	_ = http.ListenAndServe(":8080", nil)
}

// index список постов
func index(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write(static.IndexPage)
}

// post запись поста
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

// apiPosts получение списка постов
func apiPosts(w http.ResponseWriter, r *http.Request) {
	getParams := r.URL.Query()

	quantity := 10
	quantityParamName := "quantity"
	if getParams.Has(quantityParamName) {
		getQuantity, err := strconv.Atoi(getParams.Get(quantityParamName))
		if err == nil && getQuantity > 0 && getQuantity < 50 {
			quantity = getQuantity
		}
	}

	rows, err := db.DB.Query(modelPost.SqlSelect(quantity))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var posts []*modelPost.Post
	for rows.Next() {
		post := &modelPost.Post{}
		var timeInDB string

		err = rows.Scan(&post.Id, &timeInDB, &post.Text)
		if err != nil {
			panic(err)
		}
		post.Time, err = time.Parse(time.RFC3339, timeInDB)
		if err != nil {
			post.Time = time.Time{}
		}

		posts = append(posts, post)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":    true,
		"posts": posts,
	})
}

func style(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	_, _ = w.Write([]byte(static.StylePage))
}

func favicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "image/png")
	_, _ = w.Write(static.FaviconImage)
}
