/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func loginPortal(cmd *cobra.Command, args []string) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to open .env file:", err)
		return
	}

	accountName := os.Getenv("ACCOUNT_NAME")
	password := os.Getenv("PASSWORD")

	page := rod.New().MustConnect().MustPage(`
		https://portal.nap.gsic.titech.ac.jp/GetAccess/Login?Template=userpass_key&AUTHMETHOD=UserPassword
	`)

	// 	defer page.MustClose()

	page.MustElementX("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[2]/td/input").MustInput(accountName)
	page.MustElementX("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[3]/td/input").MustInput(password)

	page.MustElementX("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[5]/td/input[1]").MustClick()
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Automatically open the portal site and log in.",
	Long:  `Automatically open the portal site and log in.`,
	Run:   loginPortal,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
