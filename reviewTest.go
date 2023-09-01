package moke_kit

import "strings"

func needReview() {
	src := []int{1, 2, 3}
	dst := make([]int, 3)
	for i, x := range src {
		dst[i] = x
	}
}

func needReview1() {
	x := true
	if x == true {
	}

	x1, y1 := "1", "2"
	if strings.Index(x1, y1) != -1 {
	}
}
