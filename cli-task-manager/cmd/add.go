package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhenyang0405/gophercises/cli-task-manager/database"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task into the task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("Failed to add task:", err.Error())
		}
		fmt.Printf("Added '%s' into the task lists.", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}