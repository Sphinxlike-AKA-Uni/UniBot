package Sweeper

import (
	"fmt"
	"time"
	"math/rand"
	//"strconv"
)



func SweepWithRNG(rows, cols, mines int, r *rand.Rand) string {
	var spots [][]byte
	// Create arrays
	for r := 0; r < rows; r++ {
		s := []byte{}
		for c := 0; c < cols; c++ {
			s = append(s, 0)
		}
		spots = append(spots, s)
	}
	
	// Create mines
	m := 0
	for m < mines {
		sr := r.Intn(rows)
		sc := r.Intn(cols)
		if spots[sr][sc] == 9 {
			continue
		} else {
			spots[sr][sc] = 9
			m += 1
		}
	}
	
	// Check for neighbours
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if spots[r][c] == 9 {
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						cr := r-1+i
						cc := c-1+j
						if cc >= 0 && cc < cols {
							if cr >= 0 && cr < rows {
								if spots[cr][cc] != 9 {
									spots[cr][cc] = spots[cr][cc] + 1
								}
							}
						}
					}
				}
			}
		}
	}
	
	// Output
	rs := ""
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			rs = fmt.Sprintf("%s||%s|| ", rs, []string{":zero:", ":one:", ":two:", ":three:", ":four:", ":five:", ":six:", ":seven:", ":eight:", ":bomb:"}[spots[r][c]])
		}
		rs = fmt.Sprintf("%s\n", rs)
	}
	
	return rs
}

func Sweep(rows, cols, mines int) string {
	r := rand.New(rand.NewSource(99))
	r.Seed(time.Now().UnixNano())
	return SweepWithRNG(rows, cols, mines, r)
}