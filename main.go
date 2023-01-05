package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":8080"

var dbconn *pgx.Conn

type Picture struct {
	Fingerprint string
	CreatedAt   time.Time
}

func main() {
	var err error
	dbconn, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("connected to database")

	router := http.NewServeMux()
	router.Handle("/upload", http.MaxBytesHandler(http.HandlerFunc(uploadHandler), 1024*1024*1024))
	router.HandleFunc("/pictures", picturesListingHandler)

	log.Printf("starting server on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// client request the pictures tables
// get data from the database
//the client ask for listing all the information in the database

func picturesListingHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := dbconn.Query(context.Background(), "SELECT fingerprint, created_at from pictures order by created_at desc")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	pictures := make([]Picture, 0)
	for rows.Next() {
		var pic Picture
		if err = rows.Scan(&pic.Fingerprint, &pic.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pictures = append(pictures, pic)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	if err = enc.Encode(pictures); err != nil { //encode my pictures in w
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/* each row turned into object picture, we append into a "slice"
	for i, picture := range pictures {
		fmt.Fprintf(w, "%d. fingerprint = %s\n", i+1, picture.Fingerprint)
	}

	*/
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fingerprint := fmt.Sprintf("%x", sha1.Sum(content))

	rows, err := dbconn.Query(context.Background(), "SELECT * from pictures where fingerprint=$1", fingerprint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count := 0
	for rows.Next() {
		count = count + 1 //count++
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "content already exists", http.StatusConflict)
		return
	}

	_, err = dbconn.Exec(context.Background(), "INSERT INTO pictures(id, fingerprint) VALUES ($1,$2)", uuid.New(), fingerprint)
	if err != nil {
		log.Fatal(err)
	}

}
