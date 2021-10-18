package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"webce/cmd/web/boot"
	"webce/cmd/web/boot/router"
)

var rootCmd = &cobra.Command{
	Use:   "Web Server ",
	Short: "Web Server Start with here",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		app := router.InitRouter()
		boot.Run(app)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
