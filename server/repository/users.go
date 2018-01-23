package repository

import (
	"database/sql"

	"github.com/src-d/code-annotation/server/model"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (r *Users) Create(m *model.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (id, login, username, avatar_url) VALUES (?, ?, ?, ?)",
		m.ID, m.Login, m.Username, m.AvatarURL,
	)
	return err
}

func (r *Users) Get(id int) (*model.User, error) {
	var m model.User
	err := r.db.QueryRow("SELECT id, login, username, avatar_url FROM users WHERE id=?", id).
		Scan(&m.ID, &m.Login, &m.Username, &m.AvatarURL)
	if err != nil {
		return nil, sqlE(err)
	}
	return &m, nil
}
