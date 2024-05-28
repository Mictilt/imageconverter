package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// Default output directory for compressed images
const outputDir = "./compressed"

func main() {
	app := &cli.App{
		Name: "image-compressor",
		Usage: "Compresses an image file",
		Action: func(c *cli.Context) error {
			// Prompt the user for the choice
			fmt.Println("Choose an option:")
			fmt.Println("1. Read a directory")
			fmt.Println("2. Read a specific file")
			fmt.Println("3. Compress a PDF file")
			reader := bufio.NewReader(os.Stdin)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSuffix(choice, "\n")
			choice = strings.TrimSuffix(choice, "\r")

			switch choice {
			case "1":
				// Prompt the user for the file type to compress
				fmt.Print("Enter the file type to compress (e.g., .jpg, .png): ")
				fileType, _ := reader.ReadString('\n')
				fileType = strings.TrimSuffix(fileType, "\n")
				fileType = strings.TrimSuffix(fileType, "\r")

				// Prompt the user for the directory input
				fmt.Print("Enter the directory input path: ")
				dirInput, _ := reader.ReadString('\n')
				dirInput = strings.TrimSuffix(dirInput, "\n")
				dirInput = strings.TrimSuffix(dirInput, "\r")
				dirInput = strings.Trim(dirInput, `"`)
				// Prompt the user for the quality level
				fmt.Print("Enter the quality level (0-100): ")
				qualityInput, _ := reader.ReadString('\n')
				qualityInput = strings.TrimSuffix(qualityInput, "\n")
				qualityInput = strings.TrimSuffix(qualityInput, "\r")

				// Convert the quality input to an integer
				quality, err := strconv.Atoi(qualityInput)
				if err != nil || quality < 0 || quality > 100 {
					fmt.Println("Invalid quality level. Please enter a number between 0 and 100.")
					return err
				}

				// Process all files in the directory
				err = processDirectory(fileType, dirInput, outputDir, quality)
				if err != nil {
					return err
				}
			case "2":
				// Read a specific file
				fmt.Print("Enter the path to the image file: ")
				filePath, _ := reader.ReadString('\n')
				filePath = strings.TrimSpace(filePath)
				filePath = strings.Trim(filePath, `"`)

				// Prompt the user for the quality level
				fmt.Print("Enter the quality level (0-100): ")
				qualityInput, _ := reader.ReadString('\n')
				qualityInput = strings.TrimSuffix(qualityInput, "\n")
				qualityInput = strings.TrimSuffix(qualityInput, "\r")

				// Convert the quality input to an integer
				quality, err := strconv.Atoi(qualityInput)
				if err != nil || quality < 0 || quality > 100 {
					fmt.Println("Invalid quality level. Please enter a number between 0 and 100.")
					return err
				}

				// Ensure the output directory exists
				if err := createFolder(outputDir); err != nil {
					return err
				}

				// Open the specified file
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()
				// Read the file content
				buffer, err := io.ReadAll(file)
				if err != nil {
					return err
				}
				// Call the imageProcessing function with the quality level
				filename, err := imageProcessing(buffer, quality, outputDir, filePath)
				if err != nil {
					return err
				}

				fmt.Println("Compressed image saved as:", filename)
			case "3":
				fmt.Print("Enter the path to the PDF file: ")
				filePath, _ := reader.ReadString('\n')
				filePath = strings.TrimSpace(filePath)
				filePath = strings.Trim(filePath, `"`)

				// Prompt the user for the compression level (0-100)
				

				// Ensure the output directory exists
				if err := createFolder(outputDir); err != nil {
					return err
				}

				start := time.Now()
				outputFilePath := filepath.Join(outputDir, "compressed_"+filepath.Base(filePath))

				// Compress the PDF
				compressPDF(filePath, outputFilePath)
				

				// Get input and output file stats
				inputFileInfo, err := os.Stat(filePath)
				if err != nil {
					log.Fatalf("Fail: %v\n", err)
				}

				outputFileInfo, err := os.Stat(outputFilePath)
				if err != nil {
					log.Fatalf("Fail: %v\n", err)
				}

				// Print basic optimization statistics
				inputSize := inputFileInfo.Size()
				outputSize := outputFileInfo.Size()
				ratio := 100.0 - (float64(outputSize) / float64(inputSize) * 100.0)
				duration := float64(time.Since(start)) / float64(time.Millisecond)

				fmt.Printf("Original file: %s\n", filePath)
				fmt.Printf("Original size: %d bytes\n", inputSize)
				fmt.Printf("Optimized file: %s\n", outputFilePath)
				fmt.Printf("Optimized size: %d bytes\n", outputSize)
				fmt.Printf("Compression ratio: %.2f%%\n", ratio)
				fmt.Printf("Processing time: %.2f ms\n", duration)
			default:
				fmt.Println("Invalid choice. Please enter 1 or 2.")
			}

			fmt.Println("Press Enter to exit...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return nil
		},
	}

	// Run the CLI application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
