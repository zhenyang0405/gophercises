package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/zhenyang0405/gophercises/cli-task-manager/cmd"
	db "github.com/zhenyang0405/gophercises/cli-task-manager/database"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "cli-task-manager.db")
	err := db.Init(dbPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cmd.RootCmd.Execute()
}