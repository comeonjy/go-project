package main

import "fmt"

func main() {
	arr := []int{6, 1, 2, 7, 9,6, 1, 2, 7, 9, 3, 4, 5, 10, 8, 3, 4, 5, 10, 8}
	//arr := []int{2,1,3}
	Sort(arr)
}

func Sort(arr []int) {
	if len(arr)<=1 {
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
			}
		} else {
			j--
		}
	}
	arr[0], arr[i] = arr[i], arr[0]

	Sort(arr[:i])
	Sort(arr[i+1:])

	fmt.Println(arr)
}
