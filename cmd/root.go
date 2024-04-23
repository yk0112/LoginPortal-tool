package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Login Portal tool",
	Short: "A tool to automatically log in to the Tokyo Tech portal.",
	Long:  `A tool to automatically log in to the Tokyo Tech portal.You first need to set the credentials with the init command. Then, by executing the login command, the portal site will automatically open and you can log in.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
