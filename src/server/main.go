package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func main() {
	db, err := sql.Open("sqlite3", "./mydatabase.db")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, name TEXT NOT NULL, data BLOB NOT NULL)")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on http://localhost:8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	handleRequests()
	http.ListenAndServe(":8080", c.Handler(http.DefaultServeMux))
}

func handleRequests() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
}

func upload(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./mydatabase.db")

	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	fmt.Println(r.FormValue("file"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data := make([]byte, fileHeader.Size)
	_, err = file.Read(data)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO files (name, data) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = stmt.Exec(fileHeader.Filename, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("File uploaded successfully")
}

func download(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./mydatabase.db")

	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT name, data FROM files WHERE name = ?", id)
	var name string
	var data []byte
	err = row.Scan(&name, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))

	_, err = w.Write([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
