package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
        "os"
        "log"

	_ "github.com/mattn/go-sqlite3"
)

func setupLogger() {
	logger = log.New(os.Stdout, "testLogger: ", log.LstdFlags)
}

func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("cannot open in-memory db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("cannot create table: %v", err)
	}
}

func setupTemplates(t *testing.T) {
	var err error
	tmplBlog, err = template.New("blog").Parse(`<html>{{range .}}<h2>{{.Title}}</h2>{{end}}</html>`)
	if err != nil {
		t.Fatalf("cannot parse blog template: %v", err)
	}
	tmplAdmin, err = template.New("admin").Parse(`<html><form></form></html>`)
	if err != nil {
		t.Fatalf("cannot parse admin template: %v", err)
	}
}

func TestBlogHandler_Empty(t *testing.T) {
	setupTestDB(t)
	setupTemplates(t)

	req := httptest.NewRequest(http.MethodGet, "/blog", nil)
	w := httptest.NewRecorder()

	blogHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAdminHandler_AddPost(t *testing.T) {
        setupLogger()
	setupTestDB(t)
	setupTemplates(t)

	form := strings.NewReader("title=TestPost&content=MyContent")
	req := httptest.NewRequest(http.MethodPost, "/admin", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	adminHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303 SeeOther, got %d", resp.StatusCode)
	}

	// проверим что данные в базе
	row := db.QueryRow("SELECT title,content FROM posts LIMIT 1")
	var title, content string
	if err := row.Scan(&title, &content); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if title != "TestPost" || content != "MyContent" {
		t.Fatalf("unexpected values: got %s %s", title, content)
	}
}

func TestAdminHandler_GetForm(t *testing.T) {
	setupTestDB(t)
	setupTemplates(t)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	adminHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

