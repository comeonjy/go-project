package test

import (
	"math/rand"
	"testing"
)

func init()  {
	rand.Seed(1e9)
}

func BenchmarkSort(b *testing.B) {
	var arr []int
	for i:=0;i<30;i++{
		arr=append(arr,rand.Intn(1000))
	}
	for i := 0; i < b.N; i++ {
		insertSort(arr)
	}
	b.Log(arr)
}

func insertSort(arr []int) []int {
	var arrCopy []int
	arrCopy=append(arrCopy,arr...)

	for k1,v1:=range arrCopy {
		for k2 := range arr[:k1] {
			if arr[k1] < arr[k2] {
				arr = append(arr[:k1], arr[k1+1:]...)
				var temp []int
				temp=append(temp,arr[:k2]...)
				arr = append(append(temp,v1), arr[k2:]...)
			}
		}
	}
	return arr
}

func quickSort(arr []int) {
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

	quickSort(arr[0:i])
	quickSort(arr[i+1:])
}
