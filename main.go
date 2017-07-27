package main

import (
	"github.com/spf13/cobra"
	"liulishuo/somechat/cmd"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"fmt"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   	"somechat webapp",
	Run: 	func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	RootCmd.AddCommand(cmd.WebRunCommand)

	cobra.OnInitialize(func() {
		log.Log().Println(logs.Info("cobra initialize..."))
	})

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
