package test

import (
	"testing"
)

func BenchmarkSort(b *testing.B) {
	arr := []int{6, 1, 123, 2, 7, 9, 34, 6, 15, 2, 24, 71, 932, 34, 6, 1, 2, 73, 9, 34, 62, 1, 2, 1, 2, 7, 9}
	for i := 0; i < b.N; i++ {
		Sort(arr)
	}
	b.Log(arr)
}

func Sort(arr []int) {
	if len(arr) <= 1 {
		return
	}
	obj := arr[0]
	i, j := 0, len(arr)-1
	for i < j {
		if arr[j] < obj {
			if arr[i] > obj {
				arr[i], arr[j] = arr[j], arr[i]
			} else {
				i++
				continue
			}
		}
		j--

	}
	arr[0], arr[i] = arr[i], arr[0]

	Sort(arr[0:i])
	Sort(arr[i+1:])
}
