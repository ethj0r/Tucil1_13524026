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

// ParseBoard membaca file board config
// @return Board struct yg sudah terisi data dari file input
func ParseBoard(filename string) (*Board, error) {
	file, err := os.Open(filename)
	if err!=nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line!="" {
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
		
		for j, char:=range grid[i] {
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
	for row:=0; row<b.Size; row++ {
		for col:=0; col<b.Size; col++ {
			letter := b.Grid[row][col]
			regionMap[letter] = append(regionMap[letter], Cell{Row: row, Col: col})
		}
	}

	b.Regions = []Region{}
	for letter, cells := range regionMap {
		b.Regions = append(b.Regions, Region{
			Letter: letter,
			Cells:  cells,
		})
	}
}

func (b *Board) Validate() error {
	if b.Size==0 {
		return fmt.Errorf("Board is empty")
	}
	for i, row := range b.Grid {
		if len(row)!=b.Size {
			return fmt.Errorf("Row %d has length %d, expected %d", i, len(row), b.Size)
		}
	}

	if len(b.Regions) != b.Size {
		return fmt.Errorf("Invalid board: %d regions but board size is %d", len(b.Regions), b.Size)
	}

	for _, region := range b.Regions {
		if len(region.Cells)==0 {
			return fmt.Errorf("Region '%c' has no cells", region.Letter)
		}
	}

	totalCells := 0
	for _, region:=range b.Regions {
		totalCells+=len(region.Cells)
	}
	expectedCells := b.Size*b.Size
	if totalCells!=expectedCells {
		return fmt.Errorf("Invalid board: %d total cells, expected %d", totalCells, expectedCells)
	}
	
	return nil
}

func (b *Board) PrintBoard() {
	fmt.Println("Board:")
	for _, row:=range b.Grid {
		for _, cell:=range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Println()
	}
}

func (b *Board) PrintBoardWithQueens(queens []Cell) {
	queenMap:=make(map[Cell]bool)
	for _, q := range queens {
		queenMap[q] = true
	}

	fmt.Println("Solution:")
	for row:=0; row<b.Size; row++ {
		for col:=0; col<b.Size; col++ {
			cell := Cell{Row: row, Col: col}
			if queenMap[cell] {
				fmt.Print("# ")
			} else {
				fmt.Printf("%c ", b.Grid[row][col])
			}
		}
		fmt.Println()
	}
}

func (b *Board) GetRegionByLetter (letter rune) *Region {
	for i:=range b.Regions {
		if b.Regions[i].Letter==letter {
			return &b.Regions[i]
		}
	}
	return nil
}

func (b *Board) PrintRegions() {
	fmt.Printf("Total regions: %d\n", len(b.Regions))
	for _, region := range b.Regions {
		fmt.Printf("Region '%c': %d cells - ", region.Letter, len(region.Cells))
		for _, cell := range region.Cells {
			fmt.Printf("(%d,%d) ", cell.Row, cell.Col)
		}
		fmt.Println()
	}
}