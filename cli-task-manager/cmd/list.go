package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	db "github.com/zhenyang0405/gophercises/cli-task-manager/database"
)

var listRmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the task",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.ViewAllTasks()
		if err != nil {
			fmt.Println("Failed to read tasks from database:", err.Error())
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks in the list.")
			return
		}
		fmt.Println("You have the following tasks:")
		for idx, task := range tasks {
			fmt.Printf("%d. %s\n", idx+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listRmd)
}