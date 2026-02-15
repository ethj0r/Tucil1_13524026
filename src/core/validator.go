package core
import "math"

func isValid(queens []Cell) bool {
	n := len(queens)
	for i:=0; i<n; i++ {
		for j:=i+1; j<n; j++ {
			q1, q2 := queens[i], queens[j]

			//check row col diag
			if q1.Row == q2.Row return false
			if q1.Col == Q2.Col return false
			if math.Abs(float64(q1.Row-q2.Row)) == math.Abs(float64(q1.Col-q2.Col))
		}
	}
	return true
}

func validateInput(board *Board) error {
	if len(board.Regions) != board.Size {
		return fmt.Errorf("INVALID! %d regions but board is %d", len(board.Regions), board.Size)
	}
	for _, row := range board.Grid {
		if len(row) != board.Size {
			return fmt.Erorrf("Board is not square!")
		}
	}
	return nil
}