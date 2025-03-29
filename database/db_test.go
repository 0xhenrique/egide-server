package database

import (
    //"database/sql"
    "os"
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

func TestInitDB_Success(t *testing.T) {
    // Set DB_PATH to a temporary file
    os.Setenv("DB_PATH", "./test_egide.db")
    defer os.Remove("./test_egide.db")

    db, err := InitDB()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if db == nil {
        t.Fatal("Expected non-nil *sql.DB, got nil")
    }
    defer db.Close()
}

func TestInitDB_DefaultPath(t *testing.T) {
    // Unset DB_PATH to use default
    os.Unsetenv("DB_PATH")
    defer os.Remove("./egide.db")

    db, err := InitDB()
    if err != nil {
        t.Fatalf("Expected no error with default path, got %v", err)
    }
    if db == nil {
        t.Fatal("Expected non-nil *sql.DB with default path, got nil")
    }
    defer db.Close()
}

func TestInitDB_Failure(t *testing.T) {
    // Set DB_PATH to an invalid location
    os.Setenv("DB_PATH", "/invalid/path/egide.db")

    db, err := InitDB()
    if err == nil {
        t.Fatal("Expected an error, got none")
    }
    if db != nil {
        t.Fatal("Expected nil *sql.DB on failure, got non-nil")
    }
}

func TestCreateTables(t *testing.T) {
    os.Setenv("DB_PATH", "./test_egide.db")
    defer os.Remove("./test_egide.db")

    db, err := InitDB()
    if err != nil {
        t.Fatalf("Failed to initialize DB: %v", err)
    }
    defer db.Close()

    err = createTables(db)
    if err != nil {
        t.Fatalf("Failed to create tables: %v", err)
    }

    rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name IN ('users', 'sessions', 'websites')")
    if err != nil {
        t.Fatalf("Failed to query tables: %v", err)
    }
    defer rows.Close()

    tableCount := 0
    for rows.Next() {
        tableCount++
    }

    if tableCount != 3 {
        t.Fatalf("Expected 3 tables, found %d", tableCount)
    }
}
