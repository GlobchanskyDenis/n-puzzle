package solver

import (
	// "strconv"
	"fmt"
)

type PuzzleState struct {
	N     int
	Cells [][]uint
	PrevState *PuzzleState
	Actions []Swap
}

type Pos struct {
	Row int
	Col int
}

/*	Структура описывающая произведенную смену соседних значений.
**	Содержит позицию смены (в случае горизонтальной смены ряд у них одинаковый,
**	а колонка будет соответствовать той цифре что правее), и флаг - была это
**	горизонтальная или вертикальная смена  */
type Swap struct {
	Pos
	IsHor bool // True в случае если swap был по горизонтали
}

type Queue struct {
	queue []*PuzzleState
}

func (q *Queue) PushBack(p *PuzzleState) {
	q.queue = append(q.queue, p)
}

func (q Queue) len() int {
	return len(q.queue)
}

func (q *Queue) PopUp() *PuzzleState {
	if q.len() == 0 {
		return nil
	}
	var dst = q.queue[0]
	q.queue = q.queue[1:q.len()]
	return dst
}

/*	Матрица инициализируется числом, которое больше любого возможного  */
func CreateState(n int) *PuzzleState {
	var dst = &PuzzleState{N: n}
	for i := 0; i < n; i++ {
		var dstRow []uint
		for j := 0; j < n; j++ {
			dstRow = append(dstRow, uint(n*n))
		}
		dst.Cells = append(dst.Cells, dstRow)
	}
	return dst
}

func (p *PuzzleState) CreateChildState() *PuzzleState {
	var dst = CreateState(p.N)

	for nRow, row := range p.Cells {
		for nCol, val := range row {
			dst.changeCell(val, Pos{Row: nRow, Col: nCol})
		}
	}

	dst.PrevState = p
	dst.Actions = append([]Swap{}, p.Actions...)
	return dst
}

/*	Вывожу в стундартный вывод все клетки  */
func (p PuzzleState) Print() {
	for _, row := range p.Cells {
		print("[ ")
		for _, cell := range row {
			if cell < 10 {
				print(" ")
			}
			print(int(cell))
			print(" ")
		}
		println("]")
	}
	println("")
}

/*	Выводит в стандартный вывод ожиаемое действие над текущей матрицей в человекопонятном виде  */
func (p PuzzleState) PrintAction(action Swap) {
	if action.IsHor == true {
		fmt.Printf("horizontal row %d col %d (values %d and %d)\n", action.Row, action.Col, p.Cells[action.Row][action.Col - 1], p.Cells[action.Row][action.Col])
	} else {
		fmt.Printf("vertical row %d col %d (values %d and %d)\n", action.Row, action.Col, p.Cells[action.Row - 1][action.Col], p.Cells[action.Row][action.Col])
	}
}

/*	Максимальное число среди всех клеток */
func (p PuzzleState) findMaxValue() uint {
	var max uint
	for _, row := range p.Cells {
		for _, cell := range row {
			if max < cell {
				max = cell
			}
		}
	}
	return max
}

/*	Проверяет нахождение данного числа в матрице  */
func (p PuzzleState) isValueExist(value uint) bool {
	for _, row := range p.Cells {
		for _, cell := range row {
			if value == cell {
				return true
			}
		}
	}
	return false
}

func (p *PuzzleState) changeCell(val uint, pos Pos) {
	p.Cells[pos.Row][pos.Col] = val
}

/*	Заполняет буквой Г (по горизонтали слава направо, потом вертикально сверху вниз)
**	Проверка start < end проводится функцией выше  */
func (p *PuzzleState) fillForward(startVal uint, start, end Pos) uint {
	var val = startVal
	for col := start.Col; col <= end.Col; col++ {
		if val >= uint(p.N*p.N) {
			p.changeCell(0, Pos{Row: start.Row, Col: col})
		} else {
			p.changeCell(val, Pos{Row: start.Row, Col: col})
		}
		val++
	}
	for row := start.Row + 1; row <= end.Row; row++ {
		if val >= uint(p.N*p.N) {
			p.changeCell(0, Pos{Row: row, Col: end.Col})
		} else {
			p.changeCell(val, Pos{Row: row, Col: end.Col})
		}
		val++
	}
	return val
}

/*	Заполняет буквой Г (по горизонтали слава направо, потом вертикально сверху вниз)
**	Проверка start < end проводится функцией выше  */
func (p *PuzzleState) fillBackward(startVal uint, start, end Pos) uint {
	var val = startVal
	for col := end.Col; col >= start.Col; col-- {
		if val >= uint(p.N*p.N) {
			p.changeCell(0, Pos{Row: end.Row, Col: col})
		} else {
			p.changeCell(val, Pos{Row: end.Row, Col: col})
		}
		val++
	}
	for row := end.Row - 1; row >= start.Row; row-- {
		if val >= uint(p.N*p.N) {
			p.changeCell(0, Pos{Row: row, Col: start.Col})
		} else {
			p.changeCell(val, Pos{Row: row, Col: start.Col})
		}
		val++
	}
	return val
}

/*	Находит ответ. !! НЕ РЕШЕНИЕ, А СОСТОЯНИЕ К КОТОРОМУ НАДО СТРЕМИТЬСЯ  */
func (p *PuzzleState) FindEtalon() *PuzzleState {
	var dst = p.CreateChildState()
	var start = Pos{Row: 0, Col: 0}
	var end = Pos{Row: p.N - 1, Col: p.N - 1}
	var val uint = 1
	for start.Col <= end.Col && start.Row <= end.Row {
		val = dst.fillForward(val, start, end)
		start.Row++
		end.Col--
		val = dst.fillBackward(val, start, end)
		start.Col++
		end.Row--
	}
	return dst
}

func (p *PuzzleState) IsEqual(etalon *PuzzleState) bool {
	for i := 0; i < p.N; i++ {
		for j := 0; j < p.N; j++ {
			if p.Cells[i][j] != etalon.Cells[i][j] {
				return false
			}
		}
	}
	return true
}

/*	В случае Col == 1 будут поменяны местами колонки 0 и 1  */
func (p *PuzzleState) SwapHorizontal(pos Pos) {
	p.Actions = append(p.Actions, Swap{Pos: Pos{Row: pos.Row, Col: pos.Col}, IsHor: true})
	val0 := p.Cells[pos.Row][pos.Col - 1]
	val1 := p.Cells[pos.Row][pos.Col]
	p.changeCell(val0, pos)
	pos.Col--
	p.changeCell(val1, pos)
}

/*	В случае Row == 1 будут поменяны местами ячейки в рядах 0 и 1  */
func (p *PuzzleState) SwapVertical(pos Pos) {
	p.Actions = append(p.Actions, Swap{Pos: Pos{Row: pos.Row, Col: pos.Col}, IsHor: false})
	val0 := p.Cells[pos.Row - 1][pos.Col]
	val1 := p.Cells[pos.Row][pos.Col]
	p.changeCell(val0, pos)
	pos.Row--
	p.changeCell(val1, pos)
}

/*	Выдает TRUE только в случае если такие клетки существуют в матрице и одна из них - ноль  */
func (p PuzzleState) IsCanSwapHorizontal(pos Pos) bool {
	if pos.Col >= p.N || pos.Col < 1 {
		return false
	}
	if pos.Row >= p.N || pos.Row < 0 {
		return false
	}
	if p.Cells[pos.Row][pos.Col] != 0 && p.Cells[pos.Row][pos.Col - 1] != 0 {
		return false
	}
	if len(p.Actions) == 0 {
		return true
	}
	/*	В случае если последнее действие было точно таким же, то повторно его делать нет никакого смысла (оптимизация)  */
	var lastSwap = p.Actions[len(p.Actions) - 1]
	if lastSwap.IsHor == true && lastSwap.Pos == pos {
		return false
	}
	return true
}

/*	Выдает TRUE только в случае если такие клетки существуют в матрице и одна из них - ноль  */
func (p PuzzleState) IsCanSwapVertical(pos Pos) bool {
	if pos.Col >= p.N || pos.Col < 0 {
		return false
	}
	if pos.Row >= p.N || pos.Row < 1 {
		return false
	}
	if p.Cells[pos.Row][pos.Col] != 0 && p.Cells[pos.Row - 1][pos.Col] != 0 {
		return false
	}
	if len(p.Actions) == 0 {
		return true
	}
	/*	В случае если последнее действие было точно таким же, то повторно его делать нет никакого смысла (оптимизация)  */
	var lastSwap = p.Actions[len(p.Actions) - 1]
	if lastSwap.IsHor == false && lastSwap.Pos == pos {
		return false
	}
	return true
}

/*	Данная функция не защищена от невалидных изначальных значений изначального состояния
**	Поэтому предварительно должна быть проведена валидация  */
func Solve(initState *PuzzleState) *PuzzleState {
	var etalon = initState.FindEtalon()
	var q = &Queue{}
	var p = initState
	q.PushBack(initState)
	for p.IsEqual(etalon) == false {
		for row := 0; row < p.N; row++ {
			for col := 0; col < p.N; col++ {
				/*	Если какое-то действие доступно, делаю его и помещаю результат в очередь  */
				if p.IsCanSwapHorizontal(Pos{Col: col, Row: row}) == true {
					child := p.CreateChildState()
					child.SwapHorizontal(Pos{Col: col, Row: row})
					/*	Оптимизация  */
					if child.IsEqual(etalon) {
						return child
					}
					q.PushBack(child)
				}
				if p.IsCanSwapVertical(Pos{Col: col, Row: row}) == true {
					child := p.CreateChildState()
					child.SwapVertical(Pos{Col: col, Row: row})
					/*	Оптимизация  */
					if child.IsEqual(etalon) {
						return child
					}
					q.PushBack(child)
				}
			}
		}
		p = q.PopUp()
	}
	return p
}
