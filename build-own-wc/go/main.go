package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func countWords(value string) int {
	result := 0
	inWord := false

	for _, ch := range value {
		switch ch {
		case ' ', '\t', '\n', '\r':
			if inWord {
				inWord = false
				result++
			}
		default:
			inWord = true
		}
	}

	if inWord {
		result++
	}

	return result
}

func getFile(filename string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get path ", err)
		return nil, err
	}

	filepath := fmt.Sprintf("%s/%s", dir, filename)

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Failed to get file ", err)
		return nil, err
	}

	return file, nil
}

func getFilesize(filename string) (int, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get path ", err)
		return 0, err
	}

	filepath := fmt.Sprintf("%s/%s", dir, filename)

	fileCount, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Failed to read file ", err)
		return 0, err
	}

	return len(fileCount), nil
}

func main() {
	var isCount bool
	var isLine bool
	var isWordCount bool
	var isCharCount bool

	var cmdCcwc = &cobra.Command{
		Use:   "ccwc [path_to_file]",
		Short: "Do anything to the file",
		Long:  "You can use the command to show statistic of a file",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var file *os.File
			var err error

			if len(args) > 0 {
				file, err = getFile(args[0])
				if err != nil {
					log.Fatal("Failed to get file ", err)
					return
				}
			} else {
				file = os.Stdin
			}

			state, err := file.Stat()
			if err != nil {
				log.Fatal("Failed to get file size ", err)
				return
			}

			count := int(state.Size())

			defer file.Close()

			var (
				lineCount int
				wordCount int
				charCount int
			)

			if isCharCount {
				sc := bufio.NewReader(file)
				for {
					_, _, err := sc.ReadRune()
					if err != nil {
						break
					}
					charCount++
				}
			} else {
				sc := bufio.NewScanner(file)
				for sc.Scan() {
					lineCount++
					wordCount += countWords(sc.Text())
				}

				if err := sc.Err(); err != nil {
					log.Fatal("Failed to read file contents ", err)
					return
				}
			}

			if len(args) > 0 {
				switch {
				case isCount:
					fmt.Println("\t", count, args[0])
				case isLine:
					fmt.Println("\t", lineCount, args[0])
				case isWordCount:
					fmt.Println("\t", wordCount, args[0])
				case isCharCount:
					fmt.Println("\t", charCount, args[0])
				default:
					fmt.Println("\t", lineCount, wordCount, count, args[0])
				}
			} else {
				switch {
				case isCount:
					fmt.Println("\t", count)
				case isLine:
					fmt.Println("\t", lineCount)
				case isWordCount:
					fmt.Println("\t", wordCount)
				case isCharCount:
					fmt.Println("\t", charCount)
				default:
					fmt.Println("\t", lineCount, wordCount, count)
				}
			}
		},
	}

	var rootCmd = &cobra.Command{CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}}
	rootCmd.PersistentFlags().BoolVarP(&isCount, "count", "c", false, "Count the bytes in a file")
	rootCmd.PersistentFlags().BoolVarP(&isLine, "line", "l", false, "Count the line in a file")
	rootCmd.PersistentFlags().BoolVarP(&isWordCount, "word", "w", false, "Count words in a file")
	rootCmd.PersistentFlags().BoolVarP(&isCharCount, "multiline", "m", false, "Count characters in a file")
	rootCmd.AddCommand(cmdCcwc)
	rootCmd.SuggestionsMinimumDistance = 1
	rootCmd.Execute()
}
