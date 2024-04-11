package db

import (
	"database/sql"
	"os"

    _ "github.com/mattn/go-sqlite3"
)

func checkFileExists(path string) bool { 
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return true
}

func createFile(path string) (*os.File, error) {
    if checkFileExists(path) {
        return nil, nil
    }
    return os.Create(path)
}

func NewStore(kbppath string) (*sql.DB, error) {
    var err error
    _, err = createFile(kbppath + "/store.db")
    if err != nil {
        return nil, err
    }
    db, err := sql.Open("sqlite3", kbppath + "/store.db")
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS boards (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            board_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            description TEXT,
            status INTEGER NOT NULL,
            priority INTEGER NOT NULL,
            due_date TEXT,
            is_today BOOLEAN NOT NULL,
            is_archived BOOLEAN NOT NULL,
            FOREIGN KEY(board_id) REFERENCES boards(id)
        );
    `)
    if err != nil {
        return nil, err
    }

    return db, nil
}

