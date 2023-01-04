package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"net/http"
	"os"
)

const port = ":8080"

var dbconn *pgx.Conn

func main() {
	var err error
	dbconn, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("connected to database")

	router := http.NewServeMux()
	router.Handle("/upload", http.MaxBytesHandler(http.HandlerFunc(uploadHandler), 1024*1024*1024))

	log.Printf("starting server on %s", port)
	http.ListenAndServe(port, router)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fingerprint := fmt.Sprintf("%x", sha1.Sum(content))

	_, err = dbconn.Exec(context.Background(), "INSERT INTO pictures(id, fingerprint) VALUES ($1,$2)", uuid.New(), fingerprint)
	if err != nil {
		log.Fatal(err)
	}

}
