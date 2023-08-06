package postgres

import (
	"database/sql"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	"github.com/zumosik/telegram-go/lib/e"
	"github.com/zumosik/telegram-go/storage"
)

type Storage struct {
	DB *sql.DB
}

func New(databaseURL string) (Storage, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return Storage{}, err
	}

	return Storage{
		DB: db,
	}, nil
}

func (s Storage) Save(page *storage.Page) error {
	const errorMsg = "can't save page"

	query := "INSERT INTO pages (url, username) VALUES ($1, $2)"
	_, err := s.DB.Exec(query, page.URL, page.UserName)
	if err != nil {
		return e.Wrap(errorMsg, err)
	}

	return nil
}

func (s Storage) PickRandom(username string) (page *storage.Page, err error) {
	const errorMsg = "can't pick random page"

	query := "SELECT url, username FROM pages WHERE username = $1"
	rows, err := s.DB.Query(query, username)
	if err != nil {
		return nil, e.Wrap(errorMsg, err)
	}
	defer rows.Close()

	var pages []*storage.Page

	for rows.Next() {
		var p storage.Page
		if err := rows.Scan(&p.URL, &p.UserName); err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(pages) <= 0 {
		return nil, storage.ErrNoSavedPages
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(len(pages))

	return pages[n], nil

}

func (s Storage) List(username string) ([]*storage.Page, error) {
	const errorMsg = "can't get list of pages"

	query := "SELECT url, username FROM pages WHERE username = $1"
	rows, err := s.DB.Query(query, username)
	if err != nil {
		return nil, e.Wrap(errorMsg, err)
	}
	defer rows.Close()

	var pages []*storage.Page

	for rows.Next() {
		var p storage.Page
		if err := rows.Scan(&p.URL, &p.UserName); err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}

func (s Storage) Remove(p *storage.Page) error {
	const errorMsg = "can't remove page"

	query := "DELETE FROM pages WHERE url = $1 AND username = $2"
	_, err := s.DB.Exec(query, p.URL, p.UserName)
	if err != nil {
		return e.Wrap(errorMsg, err)
	}

	return nil

}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	const errorMsg = "can't check if page exists"

	query := "SELECT * FROM pages WHERE url = $1 AND username = $2"
	if err := s.DB.QueryRow(query, p.URL, p.UserName).Scan(); err != nil {
		return false, nil
	}

	return true, nil
}
