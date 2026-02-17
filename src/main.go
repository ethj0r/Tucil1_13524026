package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"Tucil1_13524026/src/core"
	"Tucil1_13524026/src/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--cli" {
		runCLI()
	} else {
		ui.LaunchGUI()
	}
}

func getFilename() string {
	fmt.Println("Available test cases:")
	files, err := os.ReadDir("test")
	if err != nil {
		fmt.Println("Could not read test directory")
	} else {
		for i, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
				fmt.Printf("  %d. %s\n", i+1, file.Name())
			}
		}
	}
	
	fmt.Println()
	fmt.Print("Enter filename (or full path): ")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	filename := strings.TrimSpace(input)
	
	if !strings.Contains(filename, "/") && !strings.Contains(filename, "\\") {
		filename = "test/" + filename
	}
	return filename
}

func saveSolution(board *core.Board, queens []core.Cell, originalFile string) {
	outFile := strings.Replace(originalFile, ".txt", "_solution.txt", 1)
	if outFile == originalFile {
		outFile = originalFile + "_solution.txt"
	}
	
	file, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	
	queenMap := make(map[core.Cell]bool)
	for _, q := range queens {
		queenMap[q] = true
	}
	
	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			cell := core.Cell{Row: row, Col: col}
			if queenMap[cell] {
				file.WriteString("#")
			} else {
				file.WriteString(string(board.Grid[row][col]))
			}
		}
		file.WriteString("\n")
	}
	
	fmt.Printf("Solution saved to: %s\n", outFile)
}