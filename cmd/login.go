package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sclevine/agouti"
	"github.com/spf13/cobra"
)

func convertKey(key string) string {
	pair := strings.Split(key[1:4], ",")
	if len(pair) != 2 {
		panic("Failed to parse Key.")
	}
	return pair[0] + "_" + pair[1]
}

func loginPortal(cmd *cobra.Command, args []string) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to open .env file:", err)
		return
	}

	accountName := os.Getenv("ACCOUNT_NAME")
	password := os.Getenv("PASSWORD")

	if accountName == "" || password == "" {
		fmt.Println("You have not set up an account name or password. Please use the init command.")
		return
	}

	driver := agouti.ChromeDriver()

	if err := driver.Start(); err != nil {
		fmt.Println("Failed to start driver:", err)
		return
	}

	page, err := driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open page:", err)
		return
	}

	if err := page.Navigate(`https://portal.titech.ac.jp/`); err != nil {
		fmt.Println("Failed to navigate:", err)
	        return
	}

	page.FindByXPath("/html/body/div/div[1]/div[2]/form[2]/input").Click()
	page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[2]/td/input").Fill(accountName)
	page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[3]/td/input").Fill(password)
	page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[5]/td/input[1]").Submit()

	authType, err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[1]/td").Text()
	if err == nil && !strings.HasPrefix(authType, "Matrix") {
		if strings.HasPrefix(authType, "Soft Token") {
			if err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[6]/td/select").Select("Matrix"); err != nil {
				fmt.Println(err)
				return
			}
			if err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[8]/td/input[1]").Submit(); err != nil {
				fmt.Println(err)
				return
			}
		} else if strings.HasPrefix(authType, "One-Time") {
			if err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[5]/td/select").Select("Matrix"); err != nil {
				fmt.Println(err)
				return
			}
			if err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[7]/td/input[1]").Submit(); err != nil {
				fmt.Println(err)
				return
			}
		}
		authType, _ = page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[1]/td").Text()
		if !strings.HasPrefix(authType, "Matrix") {
			fmt.Println("Failed to move to Matrix authentication page.")
			return
		}
	}

	key1, err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[5]/th[1]").Text()
	key2, err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[6]/th[1]").Text()
	key3, err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[7]/th[1]").Text()

	value1 := os.Getenv(convertKey(key1))
	value2 := os.Getenv(convertKey(key2))
	value3 := os.Getenv(convertKey(key3))

	if value1 == "" || value2 == "" || value3 == "" {
		fmt.Println("You have not set up a matrix table. Please use the init command.")
		return
	}

	page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[5]/td/input").Fill(value1)
	page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[6]/td/input").Fill(value2)
	page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[7]/td/input").Fill(value3)

	if err := page.FindByXPath("/html/body/center[3]/form/table/tbody/tr/td/table/tbody/tr[10]/td/input[1]").Submit(); err != nil {
		fmt.Println(err)
		return
	}
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
