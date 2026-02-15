package core
import
(
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Cell struct
{
	Row int
	Col int
}

type Region struct
{
	Letter rune
	Cells []Cell
}

type Board struct
{
	Size int
	Grid [][]rune
	Regions []Region
}

// parseBoard membaca file board config
// @return Board struct yg sudah terisi data dari file input
func parseBoard(filename string) (*Board, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("File is empty")
	}

	size := len(lines)
	for i, line := range lines {
		if len(line) != size {
			return nil, fmt.Errorf("Invalid board: line %d has length %d, expected %d (BOARD MUST BE SQUARE!)", i+1, len(line), size)
		}
	}

	grid := make([][]rune, size)
	for i, line := range lines {
		grid[i] = []rune(line)
		
		for j, char := range grid[i] {
			if !unicode.IsLetter(char) {
				return nil, fmt.Errorf("Invalid character '%c' at position (%d, %d): only letters allowed", char, i, j)
			}
			grid[i][j] = unicode.ToUpper(char)
		}
	}

	board := &Board{
		Size: size,
		Grid: grid,
	}
	board.extractRegions()

	return board, nil
}

func (b *Board) extractRegions() {
	regionMap := make(map[rune][]Cell)
	for row := 0; row < b.Size; row++ {
		for col := 0; col < b.Size; col++ {
			letter := b.Grid[row][col]
			regionMap[letter] = append(regionMap[letter], Cell{row, col})
		}
	}

	for letter, cells := range regionMap {
		b.Regions = append(b.Regions, Region{
			Letter: letter,
			Cells:  cells,
		})
	}
}
