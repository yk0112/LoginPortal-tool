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

func loginRun(cmd *cobra.Command, args []string) {
	driver := agouti.ChromeDriver()

	if err := driver.Start(); err != nil {
		fmt.Println("Failed to start driver:", err)
		return
	}
	loginPortal(driver)
}

func loginPortal(driver *agouti.WebDriver) {
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

	page, err := driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open page:", err)
		return
	}

	if err := page.Navigate(`https://portal.titech.ac.jp/`); err != nil {
		fmt.Println("Failed to navigate:", err)
		return
	}

	page.FindByID("portal-form").FindByXPath("form[2]/input").Click()
	page.FindByName("usr_name").Fill(accountName)
	page.FindByName("usr_password").Fill(password)
	page.FindByName("OK").Submit()

	authType, err := page.FindByID("authentication").FindByXPath("tbody/tr[1]/td").Text()
	if err == nil && !strings.HasPrefix(authType, "Matrix") {
		if strings.HasPrefix(authType, "Soft Token") {
			if err := page.FindByName("message4").Select("Matrix"); err != nil {
				fmt.Println(err)
				return
			}
			if err := page.FindByName("OK").Submit(); err != nil {
				fmt.Println(err)
				return
			}
		} else if strings.HasPrefix(authType, "One-Time") {
			if err := page.FindByName("message3").Select("Matrix"); err != nil {
				fmt.Println(err)
				return
			}
			if err := page.FindByName("OK").Submit(); err != nil {
				fmt.Println(err)
				return
			}
		}
		authType, _ = page.FindByID("authentication").FindByXPath("tbody/tr[1]/td").Text()
		if !strings.HasPrefix(authType, "Matrix") {
			fmt.Println("Failed to move to Matrix authentication page.")
			return
		}
	}

	key1, err := page.FindByID("authentication").FindByXPath("tbody/tr[5]/th[1]").Text()
	key2, err := page.FindByID("authentication").FindByXPath("tbody/tr[6]/th[1]").Text()
	key3, err := page.FindByID("authentication").FindByXPath("tbody/tr[7]/th[1]").Text()

	value1 := os.Getenv(convertKey(key1))
	value2 := os.Getenv(convertKey(key2))
	value3 := os.Getenv(convertKey(key3))

	if value1 == "" || value2 == "" || value3 == "" {
		fmt.Println("You have not set up a matrix table. Please use the init command.")
		return
	}

	page.FindByID("authentication").FindByXPath("tbody/tr[5]/td/input").Fill(value1)
	page.FindByID("authentication").FindByXPath("tbody/tr[6]/td/input").Fill(value2)
	page.FindByID("authentication").FindByXPath("tbody/tr[7]/td/input").Fill(value3)

	if err := page.FindByName("OK").Submit(); err != nil {
		fmt.Println(err)
		return
	}
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Automatically open the portal site and log in.",
	Long:  `Automatically open the portal site and log in.`,
	Run:   loginRun,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
