package main

import (
	"fmt"
	"math/rand"
)

func init()  {
	rand.Seed(1e9)
}

func main() {

	var arr []int
	for i:=0;i<30;i++{
		arr=append(arr,rand.Intn(1000))
	}
	checkSort(arr)
	checkSort(insertSort(arr))

	quickSort(arr)
	checkSort(arr)
}

// 插入排序
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



// 快速排序
func quickSort(arr []int) {
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

	quickSort(arr[:i])
	quickSort(arr[i+1:])

	//fmt.Println(arr)
}

// 排序检查
func checkSort(arr []int){
	flag:=""
	for k:=range arr{
		if k>0 && arr[k]<arr[k-1] {
			flag=fmt.Sprintf("%v",arr[:k+1])
			break
		}
	}
	if flag!="" {
		fmt.Println("排序错误",flag)
	}else{
		fmt.Println("排序正确")
	}
	fmt.Println(arr)
}
