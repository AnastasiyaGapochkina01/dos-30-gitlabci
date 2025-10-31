package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "blog.db"

var (
	tmplBlog  *template.Template
	tmplAdmin *template.Template
	db        *sql.DB
	logger    *log.Logger
)

type Post struct {
	ID      int
	Title   string
	Content string
	Created string
}

func initDB() {
	if _, statErr := os.Stat(dbFile); os.IsNotExist(statErr) {
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatalf("DB open error: %v", err)
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		if err != nil {
			log.Fatalf("Table creation error: %v", err)
		}
	}
}


func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Endpoint: %s | Method: %s | IP: %s", r.URL.Path, r.Method, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content, created FROM posts ORDER BY created DESC")
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created); err == nil {
			posts = append(posts, p)
		}
	}
	tmplBlog.Execute(w, posts)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		if title != "" && content != "" {
			_, err := db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
			if err == nil {
				logger.Printf("New post added: %s", title)
				http.Redirect(w, r, "/blog", http.StatusSeeOther)
				return
			}
		}
	}
	tmplAdmin.Execute(w, nil)
}

func main() {
	logFile, err := os.OpenFile("blog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Log file error: %v", err)
	}
	defer logFile.Close()
	logger = log.New(logFile, "blogApp: ", log.LstdFlags)

	initDB()

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		logger.Fatalf("DB connect error: %v", err)
	}
	defer db.Close()

	tmplBlog = template.Must(template.ParseFiles("templates/blog.html"))
	tmplAdmin = template.Must(template.ParseFiles("templates/admin.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/blog", blogHandler)
	mux.HandleFunc("/admin", adminHandler)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      logRequest(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Println("Starting server on :8000")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

