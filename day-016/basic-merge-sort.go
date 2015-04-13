//usr/bin/env go run $0 $@ ; exit

// basic merge sort

package main

import (
	"fmt"
	"math/rand"
)

func make_sequence(max int, num int) []int {
	ret := make([]int, num)
	for i := 0; i < num; i++ {
		ret[i] = rand.Intn(max)
	}
	return ret
}

func merge_sort(run []int) []int {
	if len(run) == 1 {
		return run
	}

	sorted_run1 := merge_sort(run[0 : len(run)/2])
	sorted_run2 := merge_sort(run[len(run)/2 : len(run)])
	merged_run := make([]int, len(sorted_run1)+len(sorted_run2))
	var idx, i, j int

	for idx, i, j = 0, 0, 0; i < len(sorted_run1) && j < len(sorted_run2); {
		if sorted_run1[i] < sorted_run2[j] {
			merged_run[idx] = sorted_run1[i]
			idx++
			i++
		} else {
			merged_run[idx] = sorted_run2[j]
			idx++
			j++
		}
	}
	for ; i < len(sorted_run1); i++ {
		merged_run[idx] = sorted_run1[i]
		idx++
	}
	for ; j < len(sorted_run2); j++ {
		merged_run[idx] = sorted_run2[j]
		idx++
	}

	return merged_run
}

func main() {
	input := make_sequence(100, 100)

	fmt.Println(merge_sort(input))
}
