package cmd

import (
	"fmt"
	db "github.com/zhenyang0405/gophercises/cli-task-manager/database"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var idArr []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Failed to parse args: %s", err)
			} else {
				idArr = append(idArr, id)
			}
		}

		tasksId, err := db.ViewAllTasks()
		if err != nil {
			fmt.Println("doCmd get all tasks went wrong:", err)
		}
		for _, id := range idArr {
			if id <= 0 || id > len(tasksId) {
				fmt.Println("Invalid task ID:", id)
				continue
			}
			taskId := tasksId[id - 1]
			err := db.DeleteTask(taskId.Key)
			if err != nil {
				fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", id, err)
			} else {
				fmt.Printf("Marked \"%d\" as completed.\n", id)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}