package core
import (
	"fmt"
	"math"
)

func isValid(queens []Cell) bool {
	n := len(queens)
	for i:=0; i<n; i++ {
		for j:=i+1; j<n; j++ {
			q1, q2 := queens[i], queens[j]

			//check row col diag
			if q1.Row == q2.Row {
				return false
			}
			if q1.Col == q2.Col {
				return false
			}

			rowDiff := math.Abs(float64(q1.Row - q2.Row))
			colDiff := math.Abs(float64(q1.Col - q2.Col))
			if rowDiff<=1 && colDiff<=1 {
				return false
			}
		}
	}
	return true
}

func ValidateInput(board *Board) error {
	if board == nil {
		return fmt.Errorf("board is nil")
	}
	if board.Size <= 0 {
		return fmt.Errorf("invalid board size: %d", board.Size)
	}
	if len(board.Regions) != board.Size {
		return fmt.Errorf("INVALID! %d regions but board is %d", len(board.Regions), board.Size)
	}
	for _, row := range board.Grid {
		if len(row)!=board.Size {
			return fmt.Errorf("Board is not square!")
		}
	}
	return nil
}