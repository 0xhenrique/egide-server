package main

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func connectDB() *sql.DB {
    driver := os.Getenv("DB_DRIVER")
    source := os.Getenv("DB_SOURCE")

    db, err := sql.Open(driver, source)
    if err != nil {
        log.Fatal("Skill issues: ", err)
    }
    return db
}
