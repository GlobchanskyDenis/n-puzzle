package solver

import (
	"math/rand"
	"testing"
	"time"
	"fmt"
)

func testingInit3x3() [][]uint {
	var dst [][]uint
	dst = append(dst, []uint{4, 2, 6})
	dst = append(dst, []uint{1, 0, 5})
	dst = append(dst, []uint{3, 8, 7})
	return dst
}

func testingInit3x3Simple() [][]uint {
	var dst [][]uint
	dst = append(dst, []uint{3, 4, 0})
	dst = append(dst, []uint{1, 2, 5})
	dst = append(dst, []uint{8, 7, 6})
	return dst
}

func testingInit3x3State() *PuzzleState {
	var dst = CreateState(3)
	var src = testingInit3x3()
	for nRow, row := range src {
		for nCol, val := range row {
			dst.changeCell(val, Pos{Row: nRow, Col: nCol})
		}
	}
	return dst
}

func testingInit3x3StateSimple() *PuzzleState {
	var dst = CreateState(3)
	var src = testingInit3x3Simple()
	for nRow, row := range src {
		for nCol, val := range row {
			dst.changeCell(val, Pos{Row: nRow, Col: nCol})
		}
	}
	return dst
}

func testingRandomizedState(n int) *PuzzleState {
	var dst = CreateState(n)
	var randSource = rand.NewSource(int64(time.Now().Second()))
	var random = rand.New(randSource)
	var value uint = uint(random.Uint32() % uint32(n * n))

	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			for dst.isValueExist(value) == true {
				value = uint(random.Uint32() % uint32(n * n))
			}
			dst.changeCell(value, Pos{Row: row, Col: col})
		}
	}
	return dst
}

func TestStatePrint(t *testing.T) {
	var p = testingInit3x3State()
	p.N = 3
	p.Print()
	
	
	p.SwapHorizontal(Pos{Row: 0, Col: 1}, p)
	p.SwapVertical(Pos{Row: 1, Col: 0}, p)
	if p.IsCanSwapVertical(Pos{Row: 1, Col: 1}) == false {
		t.Errorf("Не могу сделать swap по вертикали Pos{Row: 1, Col: 1}")
	}
	p.SwapVertical(Pos{Row: 1, Col: 1}, p)
	if p.IsCanSwapHorizontal(Pos{Row: 0, Col: 1}) == false {
		t.Errorf("Не могу сделать swap по горизонтали Pos{Row: 0, Col: 1}")
	}
	if p.IsCanSwapVertical(Pos{Row: 1, Col: 1}) == true {
		t.Errorf("Два идентичных swap нельзя делать один за другим Pos{Row: 1, Col: 1}")
	}
	
	p.Print()

	etalon := p.FindEtalon()
	etalon.Print()

	

	if p.IsEqual(p) == false {
		t.Errorf("Одна и та же структура при проверке эквивалентности дала результат FALSE")
	}

	if p.IsEqual(etalon) == true {
		t.Errorf("Две разные структуры при проверке эквивалентности дали результат TRUE")
	}

	/*	Переберу все варианты разрешенных ходов. По идее их должно быть только 2 - влево и вправо, а вниз нельзя так как предыдущее действие было таким же  */
	var posSlice []string
	for row := 0; row <= 3; row++ {
		for col := 0; col <= 3; col++ {
			/*	Если какое-то действие доступно, заношу его в слайс  */
			if p.IsCanSwapHorizontal(Pos{Col: col, Row: row}) == true {
				posSlice = append(posSlice, fmt.Sprintf("Горизонтальный swap %#v (%d %d)", Pos{Col: col, Row: row}, int(p.Cells[row][col - 1]), int(p.Cells[row][col])))
			}
			if p.IsCanSwapVertical(Pos{Col: col, Row: row}) == true {
				posSlice = append(posSlice, fmt.Sprintf("Вертикальный swap %#v (%d %d)", Pos{Col: col, Row: row}, int(p.Cells[row - 1][col]), int(p.Cells[row][col])))
			}
		}
	}

	if len(posSlice) != 2 {
		t.Errorf("Ожидалось %d варианта ходов, нашли %d", 3, len(posSlice))
	}
	t.Logf("%#v", posSlice)

	if !t.Failed() {
		t.Logf("Success")
	}
}

func TestQueue(t *testing.T) {
	var q = &Queue{}
	t.Logf("Сначала очередь пустая (%d)", q.len())

	q.PushBack(testingInit3x3State())
	if q.len() != 1 {
		t.Errorf("Добавил один элемент, ожидаемая длина %d, реальная длина %d", 1, q.len())
	} else {
		t.Logf("Добавил один элемент (%d)", q.len())
	}

	_ = q.PopUp()
	if q.len() != 0 {
		t.Errorf("Забрал один элемент, ожидаемая длина %d, реальная длина %d", 0, q.len())
	} else {
		t.Logf("Забрал один элемент (%d)", q.len())
	}
	
	_ = q.PopUp()
	_ = q.PopUp()
	_ = q.PopUp()

	q.PushBack(testingInit3x3State())
	q.PushBack(testingInit3x3State())
	q.PushBack(testingInit3x3State())

	if q.len() != 3 {
		t.Errorf("Добавил три элемента, ожидаемая длина %d, реальная длина %d", 3, q.len())
	} else {
		t.Logf("Добавил три элемента (%d)", q.len())
	}

	_ = q.PopUp()
	_ = q.PopUp()

	if q.len() != 1 {
		t.Errorf("Забрал два элемента, ожидаемая длина %d, реальная длина %d", 1, q.len())
	} else {
		t.Logf("Забрал два элемента (%d)", q.len())
	}
}

func TestSolver(t *testing.T) {
	/*	Решение матрицы 2*2 может привести к зацикливанию !! */
	var p = testingRandomizedState(3)
	// var p = testingInit3x3StateSimple()
	var etalon = p.FindEtalon()
	var result = Solve(p)

	if result.IsEqual(etalon) == false {
		t.Errorf("Etalon and result are not equal!!")
	} else {
		t.Logf("Solved by %d actions", len(result.Actions))
		p.Print()
		for _, action := range result.Actions {
			p.PrintAction(action)
			if action.IsHor == true {
				p.SwapHorizontal(action.Pos, etalon)
			} else {
				p.SwapVertical(action.Pos, etalon)
			}
			p.Print()
			
		}
	}
}

func TestRandomizer(t *testing.T) {
	/*	Решение матрицы 2*2 может привести к зацикливанию !! */
	testingRandomizedState(2).Print()
	testingRandomizedState(2).Print()
	testingRandomizedState(3).Print()
	testingRandomizedState(3).Print()
	testingRandomizedState(4).Print()
	testingRandomizedState(4).Print()
	testingRandomizedState(5).Print()
	testingRandomizedState(5).Print()
}
