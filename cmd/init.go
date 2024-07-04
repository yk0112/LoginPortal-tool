package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

const (
	rows    = 7
	columns = 10
)

type Matrix [rows][columns]rune

var (
	csvFile string
)

func readFromCSV(matrix *Matrix, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening CSV file %s: %v\n", filepath, err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for i := 0; i < rows; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading CSV %s: %v\n", filepath, err)
			os.Exit(1)
		}
		for j := 0; j < columns && j < len(record); j++ {
			if len(record[j]) > 0 {
				matrix[i][j] = rune(strings.ToUpper(record[j])[0])
			}
		}
	}
}

func readFromStdin(matrix *Matrix) {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	row, col := 0, 0
	for {
		printMatrixWithCursor(*matrix, row, col)
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyEnter, keyboard.KeyTab:
			col++
			if col >= columns {
				col = 0
				row++
			}
			if row >= rows {
				return
			}
		case keyboard.KeyArrowUp:
			row = (row - 1 + rows) % rows
		case keyboard.KeyArrowDown:
			row = (row + 1) % rows
		case keyboard.KeyArrowLeft:
			col = (col - 1 + columns) % columns
		case keyboard.KeyArrowRight:
			col = (col + 1) % columns
		case keyboard.KeyEsc:
			return
		default:
			if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' {
				matrix[row][col] = rune(strings.ToUpper(string(char))[0])
				col++
				if col >= columns {
					col = 0
					row++
				}
				if row >= rows {
					return
				}
			}
		}
	}
}

func printMatrixWithCursor(matrix Matrix, cursorRow, cursorCol int) {
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println("  A B C D E F G H I J")
	for i := 0; i < rows; i++ {
		fmt.Printf("%d ", i+1)
		for j := 0; j < columns; j++ {
			if i == cursorRow && j == cursorCol {
				fmt.Print("\033[7m") // Invert colors for cursor
			}
			if matrix[i][j] == 0 {
				fmt.Print("_ ")
			} else {
				fmt.Printf("%c ", matrix[i][j])
			}
			fmt.Print("\033[0m") // Reset colors
		}
		fmt.Println()
	}
	fmt.Println("\nUse arrow keys to navigate, Enter/Tab to move forward, ESC to finish")
}

func printMatrix(matrix Matrix) {
	fmt.Println("  A B C D E F G H I J")
	for i := 0; i < rows; i++ {
		fmt.Printf("%d ", i+1)
		for j := 0; j < columns; j++ {
			if matrix[i][j] == 0 {
				fmt.Print("_ ")
			} else {
				fmt.Printf("%c ", matrix[i][j])
			}
		}
		fmt.Println()
	}
}

func initCredential(cmd *cobra.Command, args []string) {
	var accountName, password string
	var matrix Matrix
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

	if csvFile != "" {
		readFromCSV(&matrix, csvFile)
	} else {
		readFromStdin(&matrix)
	}
	printMatrix(matrix)

	for i := 'A'; i <= 'J'; i++ {
		for j := 1; j <= 7; j++ {
			key := fmt.Sprintf("%c_%d", i, j)
                        value := (string)(matrix[j-1][(int)(i-'A')])
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
	initCmd.Flags().StringVarP(&csvFile, "csv", "c", "", "Path to the CSV file (optional)")
	rootCmd.AddCommand(initCmd)
}
