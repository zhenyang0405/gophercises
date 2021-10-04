package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var idx []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Failed to parse args: %s", err)
			} else {
				idx = append(idx, id)
			}
		}
		fmt.Println(idx)
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}