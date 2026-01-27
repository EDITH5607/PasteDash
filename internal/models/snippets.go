package models

import (
	"database/sql"
	"time"
)


type Snippet struct {
	ID int
	Title string
	Content string
	Create time.Time
	Expire time.Time
}

type SnippetModel struct{
	DB *sql.DB
}

func (m *SnippetModel) insert(title string, content string, expire int) (int, error){
	stmt := `INSERT INTO snippets (title, content, created, expires)
		   VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt,title,content,expire)
	if err!=nil {
		return 0,err
	}
	id,err := result.LastInsertId()
	if err!=nil{
		return 0,err
	}

	return int(id), nil
}

func (m *SnippetModel) get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) latest() ([]*Snippet, error) {
	return nil, nil
}