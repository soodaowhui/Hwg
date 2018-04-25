// 集合类型，array, slice, map等
package main

import "fmt"

func arr()  { // 定长数组
	//var arr1 [5]int
	arr1 := [...]int{1,2,3,4,5}
	fmt.Println(arr1)

	var arr2 = [3]int{3,6,7}
	fmt.Println(arr2)

	arr2d := [2][3]int{{1,2,3},{7,8,9}}
	fmt.Println(arr2d)
}

func sli()  { // 引用类型变长数组
	sli := []int{1,3,5,7,9}
	fmt.Println(sli)
	sli = append(sli, 11)
	fmt.Println(sli)
	fmt.Println(len(sli))
	fmt.Println(cap(sli))

	sli2 := sli
	sli2[0] = 777
	fmt.Println(sli,sli2)

	sli3 := make([]int, len(sli)+1) // copy的源和目标需要有一样的长度，或目标比源长度更长也可以
	copy(sli3, sli) // copy的源和目标容易搞错
	sli3[0] = 999
	fmt.Println(sli,sli3)

	sli4 := append([]int{}, sli...) // 注意...符号，append的方式是工程中更常用的方式
	fmt.Println(sli,sli4)

	sli5 := make([]int ,2, 3) // 不是二维数组，而是容量为3，实际长度2
	fmt.Println(sli5)

	sli6 := sli[2:4] // [) 左闭右开
	fmt.Println(sli6)
}

func mapAction()  {
	m := map[string]int{"ab":5,"ddi":7,"ff":0}
	fmt.Println(m)

	delete(m, "ab")
	fmt.Println(m)
}

func main()  {
	//arr()
	//sli()
	mapAction()
}