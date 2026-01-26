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
	return 0, nil
}

func (m *SnippetModel) get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) latest() ([]*Snippet, error) {
	return nil, nil
}