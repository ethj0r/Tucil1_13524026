package core
import ("time")

type SolverResult struct
{
	Solution []Cell
	Iterations int
	ExecutionTime time.Duration
	Found bool
}

type Solver struct
{
	board *Board
	iterations int
	updateCallback func([]Cell, int) //callback(curr combination, iteration count)
}

func NewSolver (board *Board) *Solver {
	return &Solver{
		board: board,
		iterations: 0,
	}
}

func (s *Solver) SetUpdateCallback(callback func ([]Cell, int)) {
	s.updateCallback = callback
}

func (s *Solver) Solve() SolverResult {
	startTime:=time.Now()
	s.iterations=0
	solution:=s.bruteForceSolve(0, []Cell{})
	return SolverResult{
		Solution: solution,
		Iterations: s.iterations,
		ExecutionTime: time.Since(startTime),
		Found: solution != nil,
	}
}

func (s *Solver) bruteForceSolve(regionIdx int, currQueens []Cell) []Cell {
	if regionIdx==len(s.board.Regions) {
		s.iterations++
		if s.updateCallback!=nil && s.iterations%50==0 {
			s.updateCallback(currQueens, s.iterations)
		}
		if isValid(currQueens) {
			return currQueens
		}
		return nil
	}

	region := s.board.Regions[regionIdx]
	for _, cell := range region.Cells {
		newQueens := make([]Cell, len(currQueens))
		copy(newQueens, currQueens)
		newQueens = append(newQueens, cell)

		res := s.bruteForceSolve(regionIdx+1, newQueens)
		if res != nil {
			return res
		}
	}
	return nil
}

func (s *Solver) GetTotalCombinations() int {
	if len(s.board.Regions)==0 {
		return 0
	}
	total := 1
	for _, region := range s.board.Regions {
		total *= len(region.Cells)
	}
	return total
}

func (s *Solver) solveWithAllComb() SolverResult {
	startTime:=time.Now()
	s.iterations=0

	allComb := s.generateAllComb()
	var solution []Cell
	found:=false

	for _, comb := range allComb {
		s.iterations++
		if s.updateCallback!=nil && s.iterations%50==0 {
			s.updateCallback(comb, s.iterations)
		}
		if isValid(comb) {
			solution = comb
			found = true
			break
		}
	}
	return SolverResult{
		Solution: solution,
		Iterations: s.iterations,
		ExecutionTime: time.Since(startTime),
		Found: found,
	}
}

func (s *Solver) generateAllComb() [][]Cell {
	regions:=s.board.Regions
	if len(regions)==0 {
		return [][]Cell{}
	}
	
	res:=[][]Cell{}
	for _, cell := range regions[0].Cells {
		res = append(res, []Cell{cell})
	}
	
	for i:=1; i<len(regions); i++ {
		res = cartesianProduct(res, regions[i].Cells)
	}
	return res
}

func cartesianProduct(existing [][]Cell, newCells []Cell) [][]Cell {
	res := [][]Cell{}
	for _, combo := range existing {
		for _, cell := range newCells {
			newCombo := make([]Cell, len(combo))
			copy(newCombo, combo)
			newCombo = append(newCombo, cell)
			res = append(res, newCombo)
		}
	}
	return res
}