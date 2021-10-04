package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listRmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the task",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	RootCmd.AddCommand(listRmd)
}