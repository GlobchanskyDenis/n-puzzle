package solver

import (
	"testing"
)

func TestMoveToNext(t *testing.T) {
	testCases := []struct{
		name        string
		startPos    Pos
		expectedPos Pos
		max         int
	}{
		{
			name: "0_0",
			startPos: Pos{Col: 0, Row: 0},
			expectedPos: Pos{Col: 1, Row: 0},
			max: 3,
		},{
			name: "1_0",
			startPos: Pos{Col: 1, Row: 0},
			expectedPos: Pos{Col: 2, Row: 0},
			max: 3,
		},{
			name: "2_0",
			startPos: Pos{Col: 2, Row: 0},
			expectedPos: Pos{Col: 2, Row: 1},
			max: 3,
		},{
			name: "2_1",
			startPos: Pos{Col: 2, Row: 1},
			expectedPos: Pos{Col: 2, Row: 2},
			max: 3,
		},{
			name: "2_2",
			startPos: Pos{Col: 2, Row: 2},
			expectedPos: Pos{Col: 1, Row: 2},
			max: 3,
		},{
			name: "1_2",
			startPos: Pos{Col: 1, Row: 2},
			expectedPos: Pos{Col: 0, Row: 2},
			max: 3,
		},{
			name: "0_2",
			startPos: Pos{Col: 0, Row: 2},
			expectedPos: Pos{Col: 0, Row: 1},
			max: 3,
		},{
			name: "0_1",
			startPos: Pos{Col: 0, Row: 1},
			expectedPos: Pos{Col: 1, Row: 1},
			max: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			_ = (&tc.startPos).moveToNext(tc.max)
			if tc.startPos != tc.expectedPos {
				t.Errorf("Fail: expected %#v got %#v", tc.expectedPos, tc.startPos)
			}
		})
	}
}