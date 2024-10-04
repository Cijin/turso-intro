package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
)

const (
	dbName       = "intro-db"
	syncInterval = time.Minute
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to find .env, using enviornment values instead")
	}

	primaryURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	if primaryURL == "" || authToken == "" {
		log.Println("Set env values before proceeding")
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get pwd, defaulting to './'. Error:", err)

		pwd = "./"
	}

	dir, err := os.MkdirTemp(pwd+"/temp", "libsql-*")
	if err != nil {
		log.Println("Unable to create temp directory. Error:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryURL, libsql.WithAuthToken(authToken), libsql.WithSyncInterval(syncInterval))
	if err != nil {
		log.Println("Unable to get replica connector. Error:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	queryUsers(db)
}
