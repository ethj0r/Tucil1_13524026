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

func runCLI() {
	filename := getFilename()
	
	board, err := core.ParseBoard(filename)
	if err != nil {
		fmt.Printf("Error parsing board: %v\n", err)
		return
	}
	
	if err := core.ValidateInput(board); err != nil {
		fmt.Printf("Invalid board: %v\n", err)
		return
	}
	
	board.PrintBoard()
	fmt.Println()
	
	solver := core.NewSolver(board)
	fmt.Println("Solving...")
	result := solver.Solve()
	
	if !result.Found {
		fmt.Println("No solution found!")
		return
	}
	
	board.PrintBoardWithQueens(result.Solution)
	fmt.Printf("\nWaktu pencarian: %d ms\n", result.ExecutionTime.Milliseconds())
	fmt.Printf("Banyak kasus yang ditinjau: %d kasus\n", result.Iterations)
	
	fmt.Print("\nApakah Anda ingin menyimpan solusi? (Ya/Tidak): ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	
	if answer == "ya" || answer == "y" {
		saveSolution(board, result.Solution, filename)
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