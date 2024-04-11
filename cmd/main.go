package main

import (
	"fmt"
	"kaban-board-plus/app"
	"kaban-board-plus/store/db"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func checkDirExists(path string) bool {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return true
}

func createDir(path string) error {
    if checkDirExists(path) {
        return nil
    }
    return os.Mkdir(path, 0755)
}

func main() {
    var err error
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting user home directory: ", err)
        os.Exit(1)
    }
    kbpPath := home + "/.kabanboardplus"
    err = createDir(kbpPath)
    if err != nil {
        fmt.Println("Error creating directory: ", err)
        os.Exit(1)
    }
    var f *os.File
    f, err = tea.LogToFile(kbpPath + "/kabanboardplus.log", "debug")
    db, err := db.NewStore(kbpPath)
    if err != nil {
        fmt.Println("Error initializating DB: ", err)
        os.Exit(1)
    }

    defer f.Close()
    defer db.Close()

    a := app.NewKabanBoardPlus(db)

    p := tea.NewProgram(a, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error starting program: ", err)
        os.Exit(1)
    }
}
