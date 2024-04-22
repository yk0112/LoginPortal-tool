package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func initCredential(cmd *cobra.Command, args []string) {
	var accountName, password string
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println("Failed to create configuration file:", err)
		return
	}
	defer file.Close()

	fmt.Println("your account name:")
	fmt.Scan(&accountName)

	_, err = file.WriteString(fmt.Sprintf("ACCOUNT_NAME=%s\n", accountName))
	if err != nil {
		fmt.Println("Failed to write account name to configuration file:", err)
		return
	}

	fmt.Println("password:")
	fmt.Scan(&password)

	_, err = file.WriteString(fmt.Sprintf("PASSWORD=%s\n", password))
	if err != nil {
		fmt.Println("Failed to write password to configuration file:", err)
		return
	}

	fmt.Println("matrix table:")
	for i := 'A'; i <= 'J'; i++ {
		for j := 1; j <= 7; j++ {
			key := fmt.Sprintf("%c_%d", i, j)
			var value string
			fmt.Println(key, ":")
			fmt.Scan(&value)
			_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
			if err != nil {
				fmt.Println("Failed to write configuration file to matrix table:", err)
				return
			}

		}
	}

	fmt.Println("Credentials set successfully")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set your Authentication information for Tokyo Tech Portal",
	Long:  `Set your Authentication information for Tokyo Tech Portal`,
	Run:   initCredential,
}

func init() {
	rootCmd.AddCommand(initCmd)
}
