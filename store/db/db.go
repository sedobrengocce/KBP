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
            status INTEGER DEFAULT 0 NOT NULL,
            priority INTEGER DEFAULT 1 NOT NULL,
            is_today BOOLEAN DEFAULT "false" NOT NULL,
            is_archived BOOLEAN DEFAULT "false" NOT NULL,
            done_date TEXT,
            FOREIGN KEY(board_id) REFERENCES boards(id)
        );
    `)
    if err != nil {
        return nil, err
    }

    return db, nil
}

