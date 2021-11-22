package solver

import (
	"fmt"
)

type PuzzleState struct {
	N     int
	Cells [][]uint
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
	Entropy uint // параметр характеризующий упорядоченность системы. Полностью упорядочено - 0
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
		fmt.Printf("horizontal row %d col %d (values %d and %d) entropy %d\n", action.Row, action.Col, p.Cells[action.Row][action.Col - 1], p.Cells[action.Row][action.Col], int(action.Entropy))
	} else {
		fmt.Printf("vertical row %d col %d (values %d and %d) entropy %d\n", action.Row, action.Col, p.Cells[action.Row - 1][action.Col], p.Cells[action.Row][action.Col], int(action.Entropy))
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
	if p == nil {
		println("State is NIL")
	}
	if etalon == nil {
		println("Etalon is NIL")
	}
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
func (p *PuzzleState) SwapHorizontal(pos Pos, etalon *PuzzleState) {
	p.Actions = append(p.Actions, Swap{Pos: Pos{Row: pos.Row, Col: pos.Col}, IsHor: true, Entropy: p.FindEntropy(etalon)})
	val0 := p.Cells[pos.Row][pos.Col - 1]
	val1 := p.Cells[pos.Row][pos.Col]
	p.changeCell(val0, pos)
	pos.Col--
	p.changeCell(val1, pos)
}

/*	В случае Row == 1 будут поменяны местами ячейки в рядах 0 и 1  */
func (p *PuzzleState) SwapVertical(pos Pos, etalon *PuzzleState) {
	p.Actions = append(p.Actions, Swap{Pos: Pos{Row: pos.Row, Col: pos.Col}, IsHor: false, Entropy: p.FindEntropy(etalon)})
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

func (p *PuzzleState) findPosOfValue(value uint) Pos {
	for nRow, row := range p.Cells {
		for nCol, val := range row {
			if val == value {
				return Pos{Row: nRow, Col: nCol}
			}
		}
	}
	/*	Вообще всегда должно находить значение, так что следующая строка никогда не должна выполниться  */
	return Pos{}
}

func findDeltaSqr(realPos, etalonPos Pos) uint {
	var deltaRowSqr, deltaColSqr int
	// if realPos.Row > etalonPos.Row {
	// 	deltaRowSqr = realPos.Row - etalonPos.Row
	// } else {
	// 	deltaRowSqr = etalonPos.Row - realPos.Row
	// }
	deltaRowSqr = (realPos.Row - etalonPos.Row) * (realPos.Row - etalonPos.Row)
	// if realPos.Col > etalonPos.Col {
	// 	deltaColSqr = realPos.Col - etalonPos.Col
	// } else {
	// 	deltaColSqr = etalonPos.Col - realPos.Col
	// }
	deltaColSqr = (realPos.Col - etalonPos.Col) * (realPos.Col - etalonPos.Col)
	return uint(deltaRowSqr + deltaColSqr)
}

func (p *PuzzleState) FindEntropy(etalon *PuzzleState) uint {
	var entropy uint
	for nRow, row := range etalon.Cells {
		for nCol, val := range row {
			entropy += findDeltaSqr(p.findPosOfValue(val), Pos{Row: nRow, Col: nCol})
		}
	}
	return entropy
}

/*	Проверяю чтобы энтропия системы уменьшалась (улучшилась ли энтропия за N * 2 + 1 ходов)  */
func (p *PuzzleState) IsEntropyCorrect() bool {
	var desiredEntropyStepsLen int = p.N + 5
	// if p.N > 5 {
	// 	desiredEntropyStepsLen = 4
	// } else {
	// 	desiredEntropyStepsLen = 8 - p.N
	// }
	if len(p.Actions) <= desiredEntropyStepsLen + 1 {
		return true
	}
	oldAction := p.Actions[len(p.Actions) - (desiredEntropyStepsLen)]
	lastAction := p.Actions[len(p.Actions) - 1]
	if oldAction.Entropy < lastAction.Entropy {
		return false
	}
	return true
}

/*	Вычисляет насколько улучшилась энтропия. Отрицательное число - ситуация ухудшилась  */
func (p *PuzzleState) GetEntropyImprovement() int {
	var desiredEntropyStepsLen int = 3
	if len(p.Actions) == 0 {
		return 0
	}
	if len(p.Actions) <= desiredEntropyStepsLen + 1 {
		return int(p.Actions[0].Entropy) - int(p.Actions[len(p.Actions) - 1].Entropy)
	}
	return int(p.Actions[len(p.Actions) - 1 - desiredEntropyStepsLen].Entropy) - int(p.Actions[len(p.Actions) - 1].Entropy)
}

/*	Функция должна оставлять только 20 самых перспективных вариантов и удалять все лишние
**	Реализовано сортировкой с последующим отбросом хвоста  */
func (q *Queue) Optimize() {
	var maxLen int = 500
	if q.len() < maxLen {
		return
	}
	for i := 0; i < q.len() - 1; i++ {
		for j := i + 1; j < q.len(); j++ {
			var state1 = q.queue[i]
			var state2 = q.queue[j]
			if state1.GetEntropyImprovement() < state2.GetEntropyImprovement() {
				tmp := q.queue[i]
				q.queue[i] = q.queue[j]
				q.queue[j] = tmp
			}
		}
	}
	println("Optimization")
	// for i := 0; i < maxLen - 1; i++ {
	// 	print(q.queue[i].GetEntropyImprovement())
	// 	print(" ")
	// }
	// println("")
	q.queue = q.queue[0:maxLen - 1]
}

/*	Данная функция не защищена от невалидных изначальных значений изначального состояния
**	Поэтому предварительно должна быть проведена валидация  */
func Solve(initState *PuzzleState) *PuzzleState {
	var etalon = initState.FindEtalon()
	var q = &Queue{}
	var p = initState
	// var maxQueueLength = 30
	var lastLen int
	q.PushBack(initState)
	for p.IsEqual(etalon) == false {
		for row := 0; row < p.N; row++ {
			for col := 0; col < p.N; col++ {
				/*	Если какое-то действие доступно, делаю его и помещаю результат в очередь  */
				if p.IsCanSwapHorizontal(Pos{Col: col, Row: row}) == true {
					child := p.CreateChildState()
					child.SwapHorizontal(Pos{Col: col, Row: row}, etalon)
					/*	Оптимизация  */
					if child.IsEqual(etalon) {
						return child
					}
					// if child.IsEntropyCorrect() {
						q.PushBack(child)
					// }
				}
				if p.IsCanSwapVertical(Pos{Col: col, Row: row}) == true {
					child := p.CreateChildState()
					child.SwapVertical(Pos{Col: col, Row: row}, etalon)
					/*	Оптимизация  */
					if child.IsEqual(etalon) {
						return child
					}
					// if child.IsEntropyCorrect() {
						q.PushBack(child)
					// }
				}
			}
		}
		if len(p.Actions) > lastLen {
			q.Optimize()
			lastLen = len(p.Actions)
			print("STEP ")
			println(lastLen)
		}
		p = q.PopUp()
		if p == nil {
			println("P is NIL after POP Up!!!!")
		}
		
	}
	return p
}

/*	Согласно теории по определенному алгоритму можно вычислить, возможно ли вообще свести
**	данное положение к решенному. Эта функция возвращает TRUE только если решение возможно
**	Алгоритм.
**	Для каждой клетки необходимо подсчитать количество других клеток после нее (по улитке) со значением менее чем у текущей
**	Если сумма всех этих значений будет четной - решение существует. Если нет - значит не решаемо! */
func (p *PuzzleState) isSolvable() bool {
	return false
}

// /*	Для текущей клетки подсчитываем количество других клеток после нее (по улитке) со значением менее чем у текущей  */
// func (p *PuzzleState) countCellsLessThanCurrent(currentPos Pos) uint {
// 	var currentValue = p.Cells[currentPos.Row][currentPos.Col]
// 	var nextPos = currentPos

// 	for 
// 	var end = Pos{Row: p.N - 1, Col: p.N - 1}
// 	var val uint = 1
// 	for start.Col <= end.Col && start.Row <= end.Row {
		
// 	}
// 	return dst
// }

/*	Меняет текущей позиции по улитке от левого верхнего угла по улитке к центру
**	TRUE будет в случае если текущее положение - центр  */
func (pos *Pos) moveToNext(max int) bool {
	if max % 2 == 0 { /*  Четное количество полей на одной стороне  */

		/*	Начинаем по часовой с верхней части. Тут направление следующего движения -- вправо
		**	Второе условие - ограничение по правому верхнему углу
		**	Третье условие - ограничение по левому верхнему углу  */
		if pos.Row < max / 2 && pos.Col + 1 < max - pos.Row && pos.Col + 2 < pos.Row {
			pos.Col++
			return false
		}
		/*	По часовой правая часть. Тут направление следующего движения -- вниз
		**	Второе условие - ограничение по правому нижнему углу
		**	Третьего условия нет - так как мы уже проверяли это в предыдущих пунктах  */
		if pos.Col >= max / 2 && pos.Row < pos.Col {
			pos.Row++
			return false
		}
		/*	По часовой нижняя часть. Тут направление следующего движения -- влево
		**	Второе условие - ограничение по левому нижнему углу
		**	Третьего условия нет - так как мы уже проверяли это в предыдущих пунктах  */
		if pos.Row >= max / 2 && pos.Col + 1 > max - pos.Row {
			pos.Col--
			return false
		}
		/*	Проверяю не является ли текущее положение центром. */
		if pos.Col == max / 2 - 1 && pos.Row == max / 2 {
			return true
		}
		/*	Все варианты кроме последнего исключены. Осталась только левая часть у которой следующее движение -- вверх  */
		pos.Row--
		return false

	} else { /*  Нечетное количество полей на одной стороне  */

		/*	Начинаем по часовой с верхней части. Тут направление следующего движения -- вправо
		**	Второе условие - ограничение по правому верхнему углу
		**	Третье условие - ограничение по левому верхнему углу  */
		if pos.Row <= max / 2 && pos.Col + 1 < max - pos.Row && pos.Col + 2 < pos.Row {
			pos.Col++
			return false
		}
		/*	По часовой правая часть. Тут направление следующего движения -- вниз
		**	Второе условие - ограничение по правому нижнему углу
		**	Третьего условия нет - так как мы уже проверяли это в предыдущих пунктах  */
		if pos.Col > max / 2 && pos.Row < pos.Col {
			pos.Row++
			return false
		}
		/*	По часовой нижняя часть. Тут направление следующего движения -- влево
		**	Второе условие - ограничение по левому нижнему углу
		**	Третьего условия нет - так как мы уже проверяли это в предыдущих пунктах  */
		if pos.Row > max / 2 && pos.Col + 1 > max - pos.Row {
			pos.Col--
			return false
		}
		/*	Проверяю не является ли текущее положение центром */
		if pos.Col == max / 2 && pos.Row == max / 2 {
			return true
		}
		/*	Все варианты кроме последнего исключены. Осталась только левая часть у которой следующее движение -- вверх  */
		pos.Row--
		return false
	}
}
