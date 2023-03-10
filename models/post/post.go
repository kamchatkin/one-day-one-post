package post

import (
	"fmt"
	"time"
)

// Post объект записи
type Post struct {
	Id   int64     `json:"id"`
	Text string    `json:"post"`
	Time time.Time `json:"created_at"`
}

const tableName = "posts"

const fieldCreatedAt = "created_at"
const fieldText = "text"

func (p *Post) SqlCreate() string {
	return fmt.Sprintf("insert into `%s` ('%s', '%s') values('%s', '%s')",
		tableName, fieldCreatedAt, fieldText, time.Now().Format(time.RFC3339), p.Text)
}

func SqlSelect(limit int) string {
	return fmt.Sprintf("select * from `%s` order by id desc limit %d", tableName, limit)
}
