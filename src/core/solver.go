package core
import ("time")

type SolverResult struct {
	Solution []Cell
	Iterations int
	ExecutionTime time.Duration
	Found bool
}

type Solve struct {
	board *Board
	iterations int
	updateCallback func([]Cell) //live gui update
}

func NewSolver (board *Board) *Solver {
	return &Solver{
		board: board,
		iterations: 0,
	}
}

//buat live updates di GUI
//logging purposes/live visualization
func (s *Solver) setUpdateCallback(callback func ([]Cell)){
	s.updateCallback = callback
}



func (s *Solver) Solve() SolverResult {
	startTime:=time.Now()
	solution:=s.bruteForceSolve(0, []Cell{})
	return SolverResult{
		Solution: solution,
		Iterations: s.iterations,
		ExecutionTime: time.Since(startTime),
		Found: solution != nil,
	}
}

func (s *Solver) bruteForceSolve(regionIdx int, currQueens []Cell) []Cell {
	s.iterations++
	if s.updateCallback != nil && s.iterations%100==0 {
		s.updateCallback(currQueens)
	}
	if regionIdx==len(s.board.Regions) {
		if isValid(currQueens) return currQueens
		return nil
	}

	region:=s.board.Regions[regionIdx]
	for _, cell:=range region.Cells {
		newQueens := append([]Cell{}, currQueens...)
		newQueens = append(newQueens, cell)
		res:=s.bruteForceSolve(regionIdx+1,newQueens)
		if res != nil return result
	}
	return nil
}