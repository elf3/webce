package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"webce/cmd/web/boot"
	"webce/cmd/web/router"
)

var rootCmd = &cobra.Command{
	Use:   "Web Server ",
	Short: "Web Server Start with here",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := router.InitRouter()
		boot.Run(app)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
