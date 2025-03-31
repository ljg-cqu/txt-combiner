package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sqweek/dialog"
)

// Function to remove empty lines from a given file
func removeEmptyLines(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", inputFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %v", inputFile, err)
	}

	file, err = os.Create(inputFile)
	if err != nil {
		return fmt.Errorf("failed to recreate file %s: %v", inputFile, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %v", inputFile, err)
		}
	}
	writer.Flush()
	return nil
}

// Merging files for Output File 1: File 1 lines first, then File 2 lines
func mergeFile1First(file1, file2, outputFile string) error {
	file1Content, err := os.Open(file1)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", file1, err)
	}
	defer file1Content.Close()

	file2Content, err := os.Open(file2)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", file2, err)
	}
	defer file2Content.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputFile, err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	scanner1 := bufio.NewScanner(file1Content)
	scanner2 := bufio.NewScanner(file2Content)

	for {
		var textFile1, textFile2 string

		if scanner1.Scan() {
			textFile1 = scanner1.Text()
		}
		if scanner2.Scan() {
			textFile2 = scanner2.Text()
		}

		if textFile1 == "" && textFile2 == "" {
			break
		}

		if textFile1 != "" {
			_, err := writer.WriteString(textFile1 + "\n")
			if err != nil {
				return fmt.Errorf("failed to write to output file: %v", err)
			}
		}
		if textFile2 != "" {
			_, err := writer.WriteString(textFile2 + "\n")
			if err != nil {
				return fmt.Errorf("failed to write to output file: %v", err)
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return fmt.Errorf("failed to write newline to output file: %v", err)
		}
	}

	writer.Flush()
	return nil
}

// Merging files for Output File 2: File 2 lines first, then File 1 lines
func mergeFile2First(file1, file2, outputFile string) error {
	return mergeFile1First(file2, file1, outputFile)
}

func main() {
	fmt.Println("Select the first input file:")
	inputFile1, err := dialog.File().
		Title("Select First Input File").
		Filter("Text Files (*.txt)", "txt").
		Load()
	if err != nil {
		fmt.Printf("Failed to select first input file: %v\n", err)
		return
	}
	if inputFile1 == "" {
		fmt.Println("No input file selected.")
		return
	}
	fmt.Printf("Selected first input file: %s\n", inputFile1)

	fmt.Println("Select the second input file:")
	inputFile2, err := dialog.File().
		Title("Select Second Input File").
		Filter("Text Files (*.txt)", "txt").
		Load()
	if err != nil {
		fmt.Printf("Failed to select second input file: %v\n", err)
		return
	}
	if inputFile2 == "" {
		fmt.Println("No input file selected.")
		return
	}
	fmt.Printf("Selected second input file: %s\n", inputFile2)

	outputFile1 := "output_file1_priority.txt" // File 1 priority
	outputFile2 := "output_file2_priority.txt" // File 2 priority

	files := []string{inputFile1, inputFile2}
	for _, file := range files {
		err := removeEmptyLines(file)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", file, err)
			return
		}
	}

	err = mergeFile1First(inputFile1, inputFile2, outputFile1)
	if err != nil {
		fmt.Printf("Error combining files into %s: %v\n", outputFile1, err)
		return
	}
	fmt.Printf("Files merged successfully into %s\n", outputFile1)

	err = mergeFile2First(inputFile1, inputFile2, outputFile2)
	if err != nil {
		fmt.Printf("Error combining files into %s: %v\n", outputFile2, err)
		return
	}
	fmt.Printf("Files merged successfully into %s\n", outputFile2)
}
