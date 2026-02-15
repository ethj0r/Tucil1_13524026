package core

type Cell struct {
	Row int
	Col int
}

type Region struct {
	Letter rune
	Cells  []Cell
}

type Board struct {
	Size    int
	Grid    [][]rune
	Regions []Region
}

// parseBoard membaca file board config
// @return Board struct yg sudah terisi data dari file input
func parseBoard(filename string) (*Board, error) {

	return nil, nil
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
