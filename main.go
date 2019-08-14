package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Entry struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Bookmark  uint      `json:"bookmark"`
	WrittenAt time.Time `json:"written_at"`
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func initDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/masuda?parseTime=true")
	if err != nil {
		log.Fatal("sql.Open: ", err)
		return nil, err
	}
	return db, nil
}

func GetAllEntry() []Entry {
	var entries []Entry
	db, err := initDB()
	if err != nil {
		log.Fatal("initDB: ", err)
	}
	defer db.Close()

	limit := 20
	rows, err := db.Query("select id, written_at, raw_content, title, bookmark from entry limit ?", limit)
	if err != nil {
		log.Fatal("select * from entry limit...: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.WrittenAt, &entry.Content, &entry.Title, &entry.Bookmark)
		if err != nil {
			log.Fatal("rows Scan: ", err)
		}
		entries = append(entries, entry)
	}
	return entries
}

func GetEntryByContent(query string) []Entry {
	var entries []Entry
	db, err := initDB()
	if err != nil {
		log.Fatal("initDB: ", err)
	}
	defer db.Close()

	limit := 20
	rows, err := db.Query("select id, written_at, raw_content, title, bookmark from entry where raw_content like ? limit ?", "%"+query+"%", limit)
	if err != nil {
		log.Fatal("select by content: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.WrittenAt, &entry.Content, &entry.Title, &entry.Bookmark)
		if err != nil {
			log.Fatal("rows Scan: ", err)
		}
		entries = append(entries, entry)
	}
	return entries
}

func GetEntryByTitle(query string) []Entry {
	var entries []Entry
	db, err := initDB()
	if err != nil {
		log.Fatal("initDB: ", err)
	}
	defer db.Close()

	limit := 20
	rows, err := db.Query("select id, written_at, raw_content, title, bookmark from entry where title like ? limit ?", "%"+query+"%", limit)
	if err != nil {
		log.Fatal("select by title: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.WrittenAt, &entry.Content, &entry.Title, &entry.Bookmark)
		if err != nil {
			log.Fatal("rows Scan: ", err)
		}
		entries = append(entries, entry)
	}
	return entries
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	entries := GetAllEntry()
	tmpl.Execute(w, entries)
}

func ContentSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("r.ParseForm(): ", err)
	}
	query := r.Form.Get("query")
	entries := GetEntryByContent(query)
	tmpl.Execute(w, entries)
}

func TitleSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("r.ParseForm(): ", err)
	}
	query := r.Form.Get("query")
	entries := GetEntryByTitle(query)
	tmpl.Execute(w, entries)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/content", ContentSearchHandler)
	http.HandleFunc("/title", TitleSearchHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
